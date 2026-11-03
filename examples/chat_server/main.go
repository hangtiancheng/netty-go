package main

import (
	"fmt"
	"net/http"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec/format"
	"github.com/yukino-js/netty-go/codec/frame"
	"github.com/yukino-js/netty-go/examples/transport/websocket"
)

var ManagerInst = NewManager()

func main() {

	// index page.
	websocket.DefaultOptions.ServeMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(indexHtml)
	})

	// child pipeline initializer.
	setupCodec := func(channel netty.Channel) {
		channel.Pipeline().
			// read websocket message
			AddLast(frame.PacketCodec(128)).
			// decode bytes to map[string]interface{}
			AddLast(format.JSONCodec(true, false)).
			// session recorder.
			AddLast(ManagerInst).
			// chat handler.
			AddLast(chatHandler{})
	}

	// setup bootstrap & startup server.
	netty.NewBootstrap(netty.WithChildInitializer(setupCodec), netty.WithTransport(websocket.New())).
		Listen("0.0.0.0:8080/chat").Sync()
}

type chatHandler struct{}

func (chatHandler) HandleActive(ctx netty.ActiveContext) {
	type wsTransport interface {
		Route() string
		Header() http.Header
	}

	if wst, ok := ctx.Channel().Transport().(wsTransport); ok {
		fmt.Printf("child connection from: %s, route: %s, Websocket-Key: %s, User-Agent: %s\n",
			ctx.Channel().RemoteAddr(), wst.Route(), wst.Header().Get("Sec-Websocket-Key"), wst.Header().Get("User-Agent"))
	}

	ctx.HandleActive()
}

func (chatHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {

	fmt.Printf("received child message from: %s, %v\n", ctx.Channel().RemoteAddr(), message)

	if cmd, ok := message.(map[string]interface{}); ok {
		cmd["id"] = ctx.Channel().ID()
	}

	ManagerInst.Broadcast(message)
}

func (chatHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Printf("child connection closed: %s %s\n", ctx.Channel().RemoteAddr(), ex.Error())
	ctx.HandleInactive(ex)
}
