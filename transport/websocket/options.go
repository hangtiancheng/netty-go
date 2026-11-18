package websocket

import (
	"compress/flate"
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yukino-js/netty-go/transport"
)

// DefaultOptions default websocket options.
var DefaultOptions = (&Options{
	OpCode:            websocket.TextMessage,
	Dialer:            *websocket.DefaultDialer,
	ServeMux:          http.DefaultServeMux,
	Backlog:           128,
	NoDelay:           true,
	CompressLevel:     flate.BestSpeed,
	CompressThreshold: 512,
}).Apply()

// Options defines the websocket transport configuration.
type Options struct {
	CertFile          string             `json:"certFile"`
	KeyFile           string             `json:"keyFile"`
	OpCode            int                `json:"opCode"` // TextMessage, BinaryMessage, or their bitwise OR.
	Routers           []string           `json:"routers"`
	CheckUTF8         bool               `json:"checkUTF8"`
	MaxFrameSize      int64              `json:"maxFrameSize"` // Maximum incoming message size, including all fragments.
	ReadBufferSize    int                `json:"readBufferSize"`
	WriteBufferSize   int                `json:"writeBufferSize"`
	Backlog           int                `json:"backlog"`
	NoDelay           bool               `json:"nodelay"`
	CompressEnabled   bool               `json:"compressEnabled"`
	CompressLevel     int                `json:"compressLevel"`
	CompressThreshold int64              `json:"compressThreshold"`
	TLS               *tls.Config        `json:"-"`
	Dialer            websocket.Dialer   `json:"-"`
	Upgrader          websocket.Upgrader `json:"-"`
	RequestHeader     http.Header        `json:"-"`
	ResponseHeader    http.Header        `json:"-"`
	ServeMux          *http.ServeMux     `json:"-"`
}

// Apply returns a configured copy so options can be reused by concurrent dials.
func (o *Options) Apply() *Options {
	options := *o
	o = &options
	o.RequestHeader = o.RequestHeader.Clone()
	o.ResponseHeader = o.ResponseHeader.Clone()
	if o.OpCode == 0 {
		o.OpCode = websocket.TextMessage
	}
	if o.ServeMux == nil {
		o.ServeMux = http.DefaultServeMux
	}
	if o.ReadBufferSize > 0 {
		o.Dialer.ReadBufferSize = o.ReadBufferSize
		o.Upgrader.ReadBufferSize = o.ReadBufferSize
	}
	if o.WriteBufferSize > 0 {
		o.Dialer.WriteBufferSize = o.WriteBufferSize
		o.Upgrader.WriteBufferSize = o.WriteBufferSize
	}
	o.Dialer.EnableCompression = o.CompressEnabled
	o.Upgrader.EnableCompression = o.CompressEnabled
	if o.TLS != nil {
		o.TLS = o.TLS.Clone()
	}
	if o.CertFile != "" && o.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(o.CertFile, o.KeyFile)
		if err != nil {
			panic(err)
		}
		if o.TLS == nil {
			o.TLS = &tls.Config{}
		}
		o.TLS.Certificates = append(o.TLS.Certificates, cert)
	}
	if o.TLS != nil {
		o.Dialer.TLSClientConfig = o.TLS
	}
	return o
}

type contextKey struct{}

// WithOptions wraps the websocket options.
func WithOptions(option *Options) transport.Option {
	return func(options *transport.Options) error {
		options.Context = context.WithValue(options.Context, contextKey{}, option.Apply())
		return nil
	}
}

// FromContext unwraps the websocket options.
func FromContext(ctx context.Context, def *Options) *Options {
	if v, ok := ctx.Value(contextKey{}).(*Options); ok {
		return v
	}
	return def
}
