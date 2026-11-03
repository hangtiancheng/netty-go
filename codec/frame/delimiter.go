package frame

import (
	"bytes"
	"fmt"
	"io"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
	"github.com/yukino-js/netty-go/utils"
)

// DelimiterCodec create delimiter codec
func DelimiterCodec(maxFrameLength int, delimiter string, stripDelimiter bool) codec.Codec {
	utils.AssertIf(maxFrameLength <= 0, "maxFrameLength must be a positive integer")
	utils.AssertIf(len(delimiter) <= 0, "delimiter must be nonempty string")
	return &delimiterCodec{
		maxFrameLength: maxFrameLength,
		delimiter:      []byte(delimiter),
		stripDelimiter: stripDelimiter,
	}
}

type delimiterCodec struct {
	maxFrameLength int
	delimiter      []byte
	stripDelimiter bool
}

func (*delimiterCodec) CodecName() string {
	return "delimiter-codec"
}

func (d *delimiterCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	// wrap to io.Reader
	reader := utils.MustToReader(message)

	readBuff := make([]byte, 0, 16)
	tempBuff := make([]byte, 1)
	for len(readBuff) < d.maxFrameLength {
		// read 1 byte
		n := utils.AssertLength(reader.Read(tempBuff[:]))

		// append to buffer
		readBuff = append(readBuff, tempBuff[:n]...)

		// check delimiter in received buffer
		if len(readBuff) >= len(d.delimiter) && bytes.Equal(d.delimiter, readBuff[len(readBuff)-len(d.delimiter):]) {

			// strip delimiter
			if d.stripDelimiter {
				readBuff = readBuff[:len(readBuff)-len(d.delimiter)]
			}

			// post message
			ctx.HandleRead(bytes.NewReader(readBuff))
			return
		}
	}

	utils.Assert(fmt.Errorf("frame length too large, readBuffLength(%d) >= maxFrameLength(%d)",
		len(readBuff), d.maxFrameLength))

}

func (d *delimiterCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	switch r := message.(type) {
	case []byte:
		ctx.HandleWrite([][]byte{
			// body
			r,
			// delimiter
			d.delimiter,
		})
	default:
		ctx.HandleWrite(io.MultiReader(
			// body
			utils.MustToReader(message),
			// delimiter
			bytes.NewReader(d.delimiter)),
		)
	}
}
