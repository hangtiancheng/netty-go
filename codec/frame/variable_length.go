package frame

import (
	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
	"github.com/yukino-js/netty-go/utils"
)

// VariableLengthCodec create maximum received length codec
func VariableLengthCodec(maxReadLength int) codec.Codec {
	utils.AssertIf(maxReadLength <= 0, "maxReadLength must be a positive integer")
	return &variableLengthCodec{maxReadLength: maxReadLength, buffer: make([]byte, maxReadLength)}
}

type variableLengthCodec struct {
	maxReadLength int // maximum received length
	buffer        []byte
}

func (*variableLengthCodec) CodecName() string {
	return "variable-length-codec"
}

func (v *variableLengthCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	reader := utils.MustToReader(message)
	var n = utils.AssertLength(reader.Read(v.buffer))
	ctx.HandleRead(v.buffer[:n])
}

func (*variableLengthCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	ctx.HandleWrite(message)
}
