package format

import (
	"strings"

	"github.com/yukino-js/netty-go/utils"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
)

// TextCodec create a text codec
func TextCodec() codec.Codec {
	return &textCodec{}
}

type textCodec struct{}

func (*textCodec) CodecName() string {
	return "text-codec"
}

func (*textCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	// read text bytes
	textBytes := utils.MustToBytes(message)

	// convert from []byte to string
	sb := strings.Builder{}
	sb.Write(textBytes)

	// post text
	ctx.HandleRead(sb.String())
}

func (*textCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	switch s := message.(type) {
	case string:
		ctx.HandleWrite(strings.NewReader(s))
	default:
		ctx.HandleWrite(message)
	}
}
