package tls

import (
	"context"
	"crypto/tls"

	"github.com/yukino-js/netty-go/transport"
)

// DefaultOptions default tls options
var DefaultOptions = &Options{
	TLS: &tls.Config{},
}

// Options to define the tls
type Options struct {
	CertFile        string      `json:"certFile"`
	KeyFile         string      `json:"keyFile"`
	ReadBufferSize  int         `json:"readBufferSize"`
	WriteBufferSize int         `json:"writeBufferSize"`
	TLS             *tls.Config `json:"-"`
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

// WithOptions to wrap the tls options
func WithOptions(option *Options) transport.Option {
	return func(options *transport.Options) error {
		options.Context = context.WithValue(options.Context, contextKey{}, option.Apply())
		return nil
	}
}

// FromContext to unwrap the tls options
func FromContext(ctx context.Context, def *Options) *Options {
	if v, ok := ctx.Value(contextKey{}).(*Options); ok {
		return v
	}
	return def
}
