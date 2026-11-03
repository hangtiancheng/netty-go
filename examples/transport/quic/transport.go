package quic

import (
	"net"

	"github.com/yukino-js/netty-go/transport"
)

type quicTransport struct {
	transport.Transport
	client bool
}

func newQuicTransport(conn net.Conn, quicOptions *Options, client bool) (*quicTransport, error) {
	return &quicTransport{
		Transport: transport.NewTransport(conn, quicOptions.ReadBufferSize, quicOptions.WriteBufferSize),
		client:    client,
	}, nil
}
