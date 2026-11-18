package kcp

import (
	"github.com/xtaci/kcp-go/v5"
	"github.com/yukino-js/netty-go/transport"
)

type kcpTransport struct {
	*kcp.UDPSession
	client bool
}

func newKcpTransport(conn *kcp.UDPSession, kcpOptions *Options, client bool) (*kcpTransport, error) {
	conn.SetStreamMode(true)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(kcpOptions.NoDelay, kcpOptions.Interval, kcpOptions.Resend, kcpOptions.NoCongestion)
	conn.SetMtu(kcpOptions.MTU)
	conn.SetWindowSize(kcpOptions.SndWnd, kcpOptions.RcvWnd)
	conn.SetACKNoDelay(kcpOptions.AckNodelay)

	if client {
		if err := conn.SetDSCP(kcpOptions.DSCP); nil != err {
			return nil, err
		}

		if err := conn.SetReadBuffer(kcpOptions.SockBuf); nil != err {
			return nil, err
		}

		if err := conn.SetWriteBuffer(kcpOptions.SockBuf); nil != err {
			return nil, err
		}
	}

	return &kcpTransport{UDPSession: conn, client: client}, nil
}

func (t *kcpTransport) Writev(buffs transport.Buffers) (int64, error) {
	n, err := t.UDPSession.WriteBuffers(buffs)
	return int64(n), err
}

func (t *kcpTransport) Flush() error {
	return nil
}

func (t *kcpTransport) RawTransport() interface{} {
	return t.UDPSession
}
