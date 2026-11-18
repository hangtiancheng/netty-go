package main

import (
	"fmt"

	"github.com/hangtiancheng/netty-go"
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
	if nil != err {
		panic(err)
	}

	select {
	case <-ch.Context().Done():
	case <-bootstrap.Context().Done():
	}

	fmt.Println("exited")
}
