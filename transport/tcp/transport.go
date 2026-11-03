package tcp

import (
	"net"

	"github.com/yukino-js/netty-go/transport"
)

type tcpTransport struct {
	transport.Transport
	client bool
}

func newTcpTransport(conn *net.TCPConn, tcpOptions *Options, client bool) (*tcpTransport, error) {

	if err := conn.SetKeepAlive(tcpOptions.KeepAlive); nil != err {
		return nil, err
	}

	if err := conn.SetKeepAlivePeriod(tcpOptions.KeepAlivePeriod); nil != err {
		return nil, err
	}

	if err := conn.SetLinger(tcpOptions.Linger); nil != err {
		return nil, err
	}

	if err := conn.SetNoDelay(tcpOptions.NoDelay); nil != err {
		return nil, err
	}

	if tcpOptions.SockBuf > 0 {
		if err := conn.SetReadBuffer(tcpOptions.SockBuf); nil != err {
			return nil, err
		}

		if err := conn.SetWriteBuffer(tcpOptions.SockBuf); nil != err {
			return nil, err
		}
	}

	return &tcpTransport{
		Transport: transport.NewTransport(conn, tcpOptions.ReadBufferSize, tcpOptions.WriteBufferSize),
		client:    client,
	}, nil
}
