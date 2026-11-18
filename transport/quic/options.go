package quic

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
	"github.com/yukino-js/netty-go/transport"
)

// DefaultOptions default quic options
var DefaultOptions = &Options{}

// Options to define the quic
type Options struct {
	CertFile        string       `json:"certFile"`
	KeyFile         string       `json:"keyFile"`
	ReadBufferSize  int          `json:"readBufferSize"`
	WriteBufferSize int          `json:"writeBufferSize"`
	Config          *quic.Config `json:"-"`
	TLS             *tls.Config  `json:"-"`
}

func (o *Options) Apply() *Options {
	if nil == o.TLS {
		o.TLS = &tls.Config{}
	}

	if "" != o.CertFile && "" != o.KeyFile {
		if cer, err := tls.LoadX509KeyPair(o.CertFile, o.KeyFile); nil != err {
			panic(err)
		} else {
			o.TLS.Certificates = []tls.Certificate{cer}
		}
	}

	return o
}

type contextKey struct{}

// WithOptions to wrap the quic options
func WithOptions(option *Options) transport.Option {
	return func(options *transport.Options) error {
		options.Context = context.WithValue(options.Context, contextKey{}, option.Apply())
		return nil
	}
}

// FromContext to unwrap the quic options
func FromContext(ctx context.Context, def *Options) *Options {
	if v, ok := ctx.Value(contextKey{}).(*Options); ok {
		return v
	}
	return def
}
