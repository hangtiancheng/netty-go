package tcp

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/yukino-js/netty-go/transport"
)

// New tcp factory
func New() transport.Factory {
	return new(tcpFactory)
}

type tcpFactory struct{}

func (*tcpFactory) Schemes() transport.Schemes {
	return transport.Schemes{"tcp", "tcp4", "tcp6"}
}

func (f *tcpFactory) Connect(options *transport.Options) (transport.Transport, error) {

	if err := f.Schemes().FixScheme(options.Address); nil != err {
		return nil, err
	}

	tcpOptions := FromContext(options.Context, DefaultOption)

	var d = net.Dialer{Timeout: tcpOptions.Timeout}
	conn, err := d.DialContext(options.Context, options.Address.Scheme, options.Address.Host)
	if nil != err {
		return nil, err
	}

	tt, err := newTcpTransport(conn.(*net.TCPConn), tcpOptions, true)
	if nil != err {
		_ = conn.Close()
		return nil, err
	}
	return tt, nil
}

func (f *tcpFactory) Listen(options *transport.Options) (transport.Acceptor, error) {

	if err := f.Schemes().FixScheme(options.Address); nil != err {
		return nil, err
	}

	l, err := net.Listen(options.Address.Scheme, options.AddressWithoutHost())
	if nil != err {
		return nil, err
	}

	return &tcpAcceptor{listener: l.(*net.TCPListener), options: FromContext(options.Context, DefaultOption)}, nil
}

type tcpAcceptor struct {
	listener *net.TCPListener
	options  *Options
	closed   int32
}

func (t *tcpAcceptor) Accept() (transport.Transport, error) {

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := t.listener.AcceptTCP()
		if nil != err {
			if 0 != atomic.LoadInt32(&t.closed) {
				return nil, err // listener closed
			}

			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return nil, err
		}

		// reset delay time
		tempDelay = 0

		tt, err := newTcpTransport(conn, t.options, false)
		if nil != err {
			_ = conn.Close()
			return nil, err
		}
		return tt, nil
	}
}

func (t *tcpAcceptor) Close() error {
	if atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return t.listener.Close()
	}
	return nil
}
