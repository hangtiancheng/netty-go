package transport

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/yukino-js/netty-go/utils"
)

// Option defines option function
type Option func(options *Options) error

// Options for transport
type Options struct {
	// In server side: listen address.
	// In client side: connect address.
	Address *url.URL

	// Other configure pass by context.WithValue
	Context context.Context

	// Attachment defines the object or data associated with the Channel
	Attachment interface{}
}

// AddressWithoutHost convert host:port to :port
func (lo *Options) AddressWithoutHost() string {
	_, port, err := net.SplitHostPort(lo.Address.Host)
	utils.Assert(err)
	return net.JoinHostPort("", port)
}

// Apply options
func (lo *Options) Apply(options ...Option) error {
	for _, option := range options {
		if err := option(lo); nil != err {
			return err
		}
	}
	return nil
}

// ParseOptions parse options
func ParseOptions(ctx context.Context, url string, options ...Option) (*Options, error) {
	option := &Options{Context: ctx}
	return option, option.Apply(append([]Option{withAddress(url)}, options...)...)
}

// withAddress for server listener or client dialer
func withAddress(address string) Option {
	return func(options *Options) (err error) {
		if options.Address, err = url.Parse(address); nil != err {
			// compatible host:port
			switch {
			case strings.Contains(err.Error(), "cannot contain colon"):
				options.Address, err = url.Parse(fmt.Sprintf("//%s", address))
			case strings.Contains(err.Error(), "missing protocol scheme"):
				options.Address, err = url.Parse(fmt.Sprintf("//%s", address))
			}
		}
		// default path: /
		if options.Address != nil && "" == options.Address.Path {
			options.Address.Path = "/"
		}
		return err
	}
}

// WithContext to hold other configure pass by context.WithValue
func WithContext(ctx context.Context) Option {
	return func(options *Options) error {
		options.Context = ctx
		return nil
	}
}

// WithAttachment a data associated with the Channel
func WithAttachment(attachment interface{}) Option {
	return func(options *Options) error {
		options.Attachment = attachment
		return nil
	}
}
