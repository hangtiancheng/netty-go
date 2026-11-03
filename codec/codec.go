package codec

import "github.com/yukino-js/netty-go"

// Codec defines an CodecHandler alias
type Codec = netty.CodecHandler

// Combine to wrap InboundHandler and OutboundHandler into Codec.
func Combine(name string, inbound netty.InboundHandler, outbound netty.OutboundHandler) Codec {
	return &combineCodec{name: name, InboundHandler: inbound, OutboundHandler: outbound}
}

type combineCodec struct {
	netty.InboundHandler
	netty.OutboundHandler
	name string
}

func (c *combineCodec) CodecName() string {
	return c.name
}
