# NETTY-GO

## Introduction

netty-go is heavily inspired by [netty](https://github.com/netty/netty)

## Feature

- Extensible transport support, default support TCP, [UDP, QUIC, KCP, Websocket](https://github.com/yukino-js/netty-go/transports)
- Extensible codec support
- Based on responsibility chain model
- Zero-dependency

## Documentation

- [GoDoc](https://godoc.org/github.com/yukino-js/netty-go)

## Repository layout

This repository is a Go workspace. Every module lives in its own top-level directory:

| Module     | Directory    | Import path                                |
| ---------- | ------------ | ------------------------------------------ |
| core       | `.`          | `github.com/yukino-js/netty-go`            |
| transports | `transports` | `github.com/yukino-js/netty-go/transports` |
| websocket  | `ws`         | `github.com/yukino-js/netty-go/ws`         |
| examples   | `examples`   | `github.com/yukino-js/netty-go/examples`   |

`go.work` connects the modules for local builds and tests, so changes in one module are immediately visible to the others. Each module keeps its own `go.mod` and `go.sum`, and can be built on its own with `go build ./...`.

## Examples

- [chat_server](https://github.com/yukino-js/netty-go/examples/blob/master/chat_server/main.go)
- [file_server](https://github.com/yukino-js/netty-go/examples/blob/master/file_server/main.go)
- [tcp_server](https://github.com/yukino-js/netty-go/examples/blob/master/tcp_server/main.go)
- [redis_cli](https://github.com/yukino-js/netty-go/examples/blob/master/redis_cli/main.go)
- [netty-go-examples](https://github.com/yukino-js/netty-go/examples)

## Quick Start

```go
package main

import (
	"fmt"
	"strings"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec/format"
	"github.com/yukino-js/netty-go/codec/frame"
)

func main() {
	// child pipeline initializer
	var childInitializer = func(channel netty.Channel) {
		channel.Pipeline().
			// the maximum allowable packet length is 128 bytes，use \n to split, strip delimiter
			AddLast(frame.DelimiterCodec(128, "\n", true)).
			// convert to string
			AddLast(format.TextCodec()).
			// LoggerHandler, print connected/disconnected event and received messages
			AddLast(LoggerHandler{}).
			// UpperHandler (string to upper case)
			AddLast(UpperHandler{})
	}

	// create bootstrap & listening & accepting
	netty.NewBootstrap(netty.WithChildInitializer(childInitializer)).
		Listen(":9527").Sync()
}

type LoggerHandler struct {}

func (LoggerHandler) HandleActive(ctx netty.ActiveContext) {
	fmt.Println("netty-go:", "->", "active:", ctx.Channel().RemoteAddr())
	// write welcome message
	ctx.Write("Hello I'm " + "netty-go")
}

func (LoggerHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	fmt.Println("netty-go:", "->", "handle read:", message)
	// leave it to the next handler(UpperHandler)
	ctx.HandleRead(message)
}

func (LoggerHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Println("netty-go:", "->", "inactive:", ctx.Channel().RemoteAddr(), ex)
	// disconnected，the default processing is to close the connection
	ctx.HandleInactive(ex)
}

type UpperHandler struct {}

func (UpperHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	// text to upper case
	text := message.(string)
	upText := strings.ToUpper(text)
	// write the result to the client
	ctx.Write(text + " -> " + upText)
}
```

use <code>Netcat</code> to send message

```bash
$ echo -n -e "Hello Netty-Go\nhttps://netty-go.com\n" | nc 127.0.0.1 9527
Hello I'm netty-go
Hello Netty-Go -> HELLO NETTY-GO
https://netty-go.com -> HTTPS://NETTY-GO.COM
```
