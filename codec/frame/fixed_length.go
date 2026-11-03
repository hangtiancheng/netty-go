package frame

import (
	"io"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
	"github.com/yukino-js/netty-go/utils"
)

// FixedLengthCodec create fixed length codec
func FixedLengthCodec(length int) codec.Codec {
	utils.AssertIf(length <= 0, "invalid fixed length")
	return &fixedLengthCodec{length}
}

type fixedLengthCodec struct {
	length int
}

func (*fixedLengthCodec) CodecName() string {
	return "fixed-length-codec"
}

func (f *fixedLengthCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	ctx.HandleRead(io.LimitReader(utils.MustToReader(message), int64(f.length)))
}

func (f *fixedLengthCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	// 直接交给下一个处理器处理
	ctx.HandleWrite(message)
}
