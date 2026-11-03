package main

import (
	"fmt"
	"net/http"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec/xhttp"
	"github.com/yukino-js/netty-go/examples/transport/tls"
)

func main() {

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, netty-go!"))
	})

	// child pipeline initializer.
	setupCodec := func(channel netty.Channel) {
		channel.Pipeline().
			// decode http request from channel
			AddLast(xhttp.ServerCodec()).
			// print http access log
			AddLast(new(httpStateHandler)).
			// compatible with http.Handler
			AddLast(xhttp.Handler(httpMux))
	}

	// setup bootstrap & startup server.
	netty.NewBootstrap(netty.WithChildInitializer(setupCodec), netty.WithTransport(tls.New())).
		Listen("0.0.0.0:8080", tls.WithOptions(&tls.Options{CertFile: "cert.pem", KeyFile: "key.pem"})).Sync()
}

type httpStateHandler struct{}

func (*httpStateHandler) HandleActive(ctx netty.ActiveContext) {
	fmt.Printf("http client active: %s\n", ctx.Channel().RemoteAddr())
	ctx.HandleActive()
}

func (*httpStateHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	if request, ok := message.(*http.Request); ok {
		fmt.Printf("[%d]%s: %s %s\n", ctx.Channel().ID(), ctx.Channel().RemoteAddr(), request.Method, request.URL.Path)
	}
	ctx.HandleRead(message)
}

func (*httpStateHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Printf("http client inactive: %s %v\n", ctx.Channel().RemoteAddr(), ex)
	ctx.HandleInactive(ex)
}
