package xhttp

import (
	"fmt"
	"net/http"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
	"github.com/yukino-js/netty-go/utils"
)

// ClientCodec create a http client codec
func ClientCodec() codec.Codec {
	return codec.Combine("http-client-codec", new(responseCodec), new(requestCodec))
}

// ServerCodec create a http server codec
func ServerCodec() codec.Codec {
	return codec.Combine("http-server-codec", new(requestCodec), new(responseCodec))
}

// Handler to convert http.Handler to codec.Codec
func Handler(handler http.Handler) codec.Codec {
	if nil == handler {
		handler = http.DefaultServeMux
	}
	return &handlerAdapter{handler: handler}
}

type handlerAdapter struct {
	handler http.Handler
}

func (*handlerAdapter) CodecName() string {
	return "http-handler-adapter"
}

func (h *handlerAdapter) HandleRead(ctx netty.InboundContext, message netty.Message) {

	switch r := message.(type) {
	case *http.Request:
		// create a response writer
		writer := NewResponseWriter(r, ctx.Channel().Writer())
		// flush response
		defer writer.Flush()
		// serve the http request
		h.handler.ServeHTTP(writer, r)
	default:
		utils.Assert(fmt.Errorf("unrecognized type: %T", message))
	}
}

func (*handlerAdapter) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	panic("The request handling has been handed over to http.Handler, and nothing should be written here")
}

func (*handlerAdapter) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	ctx.Close(ex)
}
