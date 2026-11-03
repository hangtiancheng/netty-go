package main

import (
	"fmt"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

func main() {

	// setup client pipeline initializer.
	setupCodec := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(&simpleRedisCodec{}, &simpleRedisConsole{})
	}

	// new bootstrap
	var bootstrap = netty.NewBootstrap(netty.WithClientInitializer(setupCodec))

	// connect to redis server
	fmt.Println("connecting redis server ...")

	ch, err := bootstrap.Connect("127.0.0.1:6379")
	utils.Assert(err)

	select {
	case <-ch.Context().Done():
	case <-bootstrap.Context().Done():
	}

	fmt.Println("exited")
}
