package frame

import (
	"bytes"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
	"github.com/yukino-js/netty-go/utils"
)

// PacketCodec create packet codec
func PacketCodec(readBuffSize int) codec.Codec {
	return &packetCodec{readBuffer: bytes.NewBuffer(make([]byte, 0, readBuffSize))}
}

type packetCodec struct {
	readBuffer *bytes.Buffer
}

func (packetCodec) CodecName() string {
	return "packet-codec"
}

func (p packetCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	// reset read buffer
	p.readBuffer.Reset()
	// read packet
	utils.AssertLong(p.readBuffer.ReadFrom(utils.MustToReader(message)))
	// post packet
	ctx.HandleRead(p.readBuffer)
}

func (packetCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	ctx.HandleWrite(message)
}
