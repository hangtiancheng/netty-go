package frame

import (
	"io"

	"github.com/hangtiancheng/netty-go"
	"github.com/hangtiancheng/netty-go/codec"
	"github.com/hangtiancheng/netty-go/internal/utils"
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
