package tcp

import (
	"context"
	"time"

	"github.com/yukino-js/netty-go/transport"
)

// DefaultOption default tcp options
var DefaultOption = &Options{
	Timeout:         time.Second * 5,
	KeepAlive:       true,
	KeepAlivePeriod: time.Minute,
	Linger:          -1,
	NoDelay:         true,
}

// Options fot tcp transport
type Options struct {
	Timeout         time.Duration `json:"timeout"`
	KeepAlive       bool          `json:"keep-alive"`
	KeepAlivePeriod time.Duration `json:"keep-alive-period"`
	Linger          int           `json:"linger"`
	NoDelay         bool          `json:"nodelay"`
	SockBuf         int           `json:"sockbuf"`
	ReadBufferSize  int           `json:"readBufferSize"`
	WriteBufferSize int           `json:"writeBufferSize"`
}

type contextKey struct{}

// WithOptions to wrap the tcp options
func WithOptions(option *Options) transport.Option {
	return func(options *transport.Options) error {
		options.Context = context.WithValue(options.Context, contextKey{}, option)
		return nil
	}
}

// FromContext to unwrap the tcp options
func FromContext(ctx context.Context, def *Options) *Options {
	if v, ok := ctx.Value(contextKey{}).(*Options); ok {
		return v
	}
	return def
}
