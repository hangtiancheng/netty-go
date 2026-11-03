package frame

import (
	"encoding/binary"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

// LengthFieldPrepender for LengthFieldCodec
//
// An encoder to prepends the length of the message
func LengthFieldPrepender(
	byteOrder binary.ByteOrder,
	lengthFieldLength int,
	lengthAdjustment int,
	lengthIncludesLengthFieldLength bool,
) netty.OutboundHandler {
	utils.AssertIf(lengthFieldLength != 1 && lengthFieldLength != 2 &&
		lengthFieldLength != 4 && lengthFieldLength != 8, "lengthFieldLength must be either 1, 2, 3, 4, or 8")
	return &lengthFieldPrepender{
		byteOrder:                       byteOrder,
		lengthFieldLength:               lengthFieldLength,
		lengthAdjustment:                lengthAdjustment,
		lengthIncludesLengthFieldLength: lengthIncludesLengthFieldLength,
	}
}

type lengthFieldPrepender struct {
	byteOrder                       binary.ByteOrder
	lengthFieldLength               int
	lengthAdjustment                int
	lengthIncludesLengthFieldLength bool
}

func (l *lengthFieldPrepender) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	bodyBytes := utils.MustToBytes(message)

	length := len(bodyBytes) + l.lengthAdjustment
	if l.lengthIncludesLengthFieldLength {
		length += l.lengthFieldLength
	}

	// head buffer
	lengthBuff := packFieldLength(l.byteOrder, l.lengthFieldLength, int64(length))

	// HEAD | BODY
	ctx.HandleWrite([][]byte{lengthBuff, bodyBytes})
}
