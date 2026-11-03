package kcp

import (
	"github.com/xtaci/kcp-go/v5"
	"github.com/yukino-js/netty-go/transport"
)

// New a kcp transport factory
func New() transport.Factory {
	return new(kcpFactory)
}

type kcpFactory struct{}

func (*kcpFactory) Schemes() transport.Schemes {
	return transport.Schemes{"kcp"}
}

func (f *kcpFactory) Connect(options *transport.Options) (transport.Transport, error) {

	if err := f.Schemes().FixScheme(options.Address); nil != err {
		return nil, err
	}

	kcpOptions := FromContext(options.Context, DefaultOptions)

	conn, err := kcp.DialWithOptions(options.Address.Host, kcpOptions.Block, kcpOptions.DataShard, kcpOptions.ParityShard)
	if nil != err {
		return nil, err
	}

	tt, err := newKcpTransport(conn, kcpOptions, true)
	if nil != err {
		_ = conn.Close()
		return nil, err
	}
	return tt, nil
}

func (f *kcpFactory) Listen(options *transport.Options) (transport.Acceptor, error) {

	if err := f.Schemes().FixScheme(options.Address); nil != err {
		return nil, err
	}

	kcpOptions := FromContext(options.Context, DefaultOptions)

	l, err := kcp.ListenWithOptions(options.AddressWithoutHost(), kcpOptions.Block, kcpOptions.DataShard, kcpOptions.ParityShard)
	if nil != err {
		return nil, err
	}

	if err = l.SetDSCP(kcpOptions.DSCP); nil != err {
		_ = l.Close()
		return nil, err
	}

	if err = l.SetReadBuffer(kcpOptions.SockBuf); nil != err {
		_ = l.Close()
		return nil, err
	}

	if err = l.SetWriteBuffer(kcpOptions.SockBuf); nil != err {
		_ = l.Close()
		return nil, err
	}

	return &kcpAcceptor{listener: l, options: kcpOptions}, nil
}

type kcpAcceptor struct {
	listener *kcp.Listener
	options  *Options
}

func (k *kcpAcceptor) Accept() (transport.Transport, error) {
	conn, err := k.listener.AcceptKCP()
	if nil != err {
		return nil, err
	}

	tt, err := newKcpTransport(conn, k.options, false)
	if nil != err {
		_ = conn.Close()
		return nil, err
	}
	return tt, nil
}

func (k *kcpAcceptor) Close() error {
	if k.listener != nil {
		defer func() { k.listener = nil }()
		return k.listener.Close()
	}
	return nil
}
