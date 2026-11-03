package netty

import (
	"context"
	"sync/atomic"

	"github.com/yukino-js/netty-go/transport"
)

type (
	// ChannelInitializer to init the pipeline of channel
	ChannelInitializer func(Channel)
	// ChannelFactory to create a channel
	ChannelFactory func(id int64, ctx context.Context, pipeline Pipeline, transport transport.Transport, executor Executor) Channel
	// PipelineFactory to create pipeline
	PipelineFactory func() Pipeline
	// TransportFactory tp create transport
	TransportFactory transport.Factory
	// ChannelIDFactory to create channel id
	ChannelIDFactory func() int64
	// ChannelHolder to manage channels
	ChannelHolder interface {
		ActiveHandler
		InactiveHandler
		// CloseAll close the all channels
		CloseAll(err error)
	}

	// bootstrapOptions
	bootstrapOptions struct {
		bootstrapCtx      context.Context
		bootstrapCancel   context.CancelFunc
		clientInitializer ChannelInitializer
		childInitializer  ChannelInitializer
		transportFactory  TransportFactory
		channelFactory    ChannelFactory
		pipelineFactory   PipelineFactory
		channelIDFactory  ChannelIDFactory
		executor          Executor
		holder            ChannelHolder
	}
)

// SequenceID to generate a sequence id
func SequenceID() ChannelIDFactory {
	var id int64
	return func() int64 {
		return atomic.AddInt64(&id, 1)
	}
}

type Option func(options *bootstrapOptions)

// WithContext fork child context with context.WithCancel
func WithContext(ctx context.Context) Option {
	return func(options *bootstrapOptions) {
		options.bootstrapCtx, options.bootstrapCancel = context.WithCancel(ctx)
	}
}

// WithChannelID to set ChannelIDFactory
func WithChannelID(channelIDFactory ChannelIDFactory) Option {
	return func(options *bootstrapOptions) {
		options.channelIDFactory = channelIDFactory
	}
}

// WithPipeline to set PipelineFactory
func WithPipeline(pipelineFactory PipelineFactory) Option {
	return func(options *bootstrapOptions) {
		options.pipelineFactory = pipelineFactory
	}
}

// WithChannel to set ChannelFactory
func WithChannel(channelFactory ChannelFactory) Option {
	return func(options *bootstrapOptions) {
		options.channelFactory = channelFactory
	}
}

// WithTransport to set TransportFactory
func WithTransport(transportFactory TransportFactory) Option {
	return func(options *bootstrapOptions) {
		options.transportFactory = transportFactory
	}
}

// WithChildInitializer to set server side ChannelInitializer
func WithChildInitializer(initializer ChannelInitializer) Option {
	return func(options *bootstrapOptions) {
		options.childInitializer = initializer
	}
}

// WithClientInitializer to set client side ChannelInitializer
func WithClientInitializer(initializer ChannelInitializer) Option {
	return func(options *bootstrapOptions) {
		options.clientInitializer = initializer
	}
}

// WithExecutor use custom Executor
func WithExecutor(executor Executor) Option {
	return func(options *bootstrapOptions) {
		options.executor = executor
	}
}

// WithChannelHolder use custom ChannelHolder
func WithChannelHolder(holder ChannelHolder) Option {
	return func(options *bootstrapOptions) {
		options.holder = holder
	}
}
