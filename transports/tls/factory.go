package tls

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/yukino-js/netty-go/transport"
)

// New a tls transport factory
func New() transport.Factory {
	return new(tlsFactory)
}

type tlsFactory struct{}

func (t *tlsFactory) Schemes() transport.Schemes {
	return transport.Schemes{"tcp", "tcp4", "tcp6"}
}

func (t *tlsFactory) Connect(options *transport.Options) (transport.Transport, error) {

	if err := t.Schemes().FixScheme(options.Address); nil != err {
		return nil, err
	}

	tlsOptions := FromContext(options.Context, DefaultOptions)

	conn, err := tls.Dial(options.Address.Scheme, options.Address.Host, tlsOptions.TLS)
	if nil != err {
		return nil, err
	}

	tt, err := newTlsTransport(conn, tlsOptions, true)
	if nil != err {
		_ = conn.Close()
		return nil, err
	}
	return tt, nil
}

func (t *tlsFactory) Listen(options *transport.Options) (transport.Acceptor, error) {

	if err := t.Schemes().FixScheme(options.Address); nil != err {
		return nil, err
	}

	tlsOptions := FromContext(options.Context, DefaultOptions)

	l, err := tls.Listen(options.Address.Scheme, options.AddressWithoutHost(), tlsOptions.TLS)
	if nil != err {
		return nil, err
	}

	return &tlsAcceptor{listener: l, options: tlsOptions}, nil
}

type tlsAcceptor struct {
	listener net.Listener
	options  *Options
}

func (t *tlsAcceptor) Accept() (transport.Transport, error) {
	if nil == t.listener {
		return nil, errors.New("no listener")
	}

	conn, err := t.listener.Accept()
	if nil != err {
		return nil, err
	}

	tt, err := newTlsTransport(conn.(*tls.Conn), t.options, false)
	if nil != err {
		_ = conn.Close()
		return nil, err
	}
	return tt, nil
}

func (t *tlsAcceptor) Close() error {
	if t.listener != nil {
		defer func() { t.listener = nil }()
		return t.listener.Close()
	}
	return nil
}
