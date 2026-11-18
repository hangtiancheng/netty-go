package websocket

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/transport"
)

type Upgrader interface {
	Upgrade(writer http.ResponseWriter, request *http.Request) (netty.Channel, error)
}

type channelServe interface {
	ServeChannel(ctx context.Context, transport transport.Transport, attachment netty.Attachment, childChannel bool) netty.Channel
}

type HTTPUpgrader struct {
	Upgrader   websocket.Upgrader
	ctx        context.Context
	attachment netty.Attachment
	options    *Options
	serve      channelServe
}

func NewHTTPUpgrader(engine netty.Bootstrap, option ...transport.Option) HTTPUpgrader {
	options, err := transport.ParseOptions(engine.Context(), "ws://127.0.0.1:0", option...)
	if nil != err {
		panic(err)
	}

	wsOptions := FromContext(options.Context, DefaultOptions)

	return HTTPUpgrader{
		Upgrader:   wsOptions.Upgrader,
		ctx:        options.Context,
		attachment: options.Attachment,
		options:    wsOptions,
		serve:      engine.(channelServe),
	}
}

func (hu HTTPUpgrader) Upgrade(writer http.ResponseWriter, request *http.Request) (netty.Channel, error) {
	conn, err := hu.Upgrader.Upgrade(writer, request, hu.options.ResponseHeader)
	if nil != err {
		if nil != conn {
			_ = conn.Close()
		}
		return nil, err
	}

	t, err := newWebsocketTransport(conn, hu.options, request)
	if nil != err {
		_ = conn.Close()
		return nil, err
	}

	return hu.serve.ServeChannel(hu.ctx, t, hu.attachment, true), nil
}
