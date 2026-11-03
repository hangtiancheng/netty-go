package tls

import (
	"crypto/tls"

	"github.com/yukino-js/netty-go/transport"
)

type tlsTransport struct {
	transport.Transport
	client bool
}

func newTlsTransport(conn *tls.Conn, tlsOptions *Options, client bool) (*tlsTransport, error) {
	return &tlsTransport{
		Transport: transport.NewTransport(conn, tlsOptions.ReadBufferSize, tlsOptions.WriteBufferSize),
		client:    client,
	}, nil
}
