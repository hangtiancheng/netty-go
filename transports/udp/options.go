package udp

import (
	"context"

	"github.com/yukino-js/netty-go/transport"
)

// DefaultOptions default udp options
var DefaultOptions = &Options{
	MaxPacketSize: 1400,
	MaxBacklog:    16,
	ReusePort:     false,
}

// Options to define the udp
type Options struct {
	MaxPacketSize int32 `json:"max-packet-size"`
	MaxBacklog    int32 `json:"max-backlog"`
	ReusePort     bool  `json:"reuse-port"`
}

type contextKey struct{}

// WithOptions to wrap the udp options
func WithOptions(option *Options) transport.Option {
	return func(options *transport.Options) error {
		options.Context = context.WithValue(options.Context, contextKey{}, option)
		return nil
	}
}

// FromContext to unwrap the udp options
func FromContext(ctx context.Context, def *Options) *Options {
	if v, ok := ctx.Value(contextKey{}).(*Options); ok {
		return v
	}
	return def
}
