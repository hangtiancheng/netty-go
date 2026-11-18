# NETTY-GO

## Introduction

netty-go is heavily inspired by [netty](https://github.com/netty/netty)

## Feature

- Extensible transport support, default support TCP, [UDP, QUIC, KCP, Websocket](./transport)
- Extensible codec support
- Based on responsibility chain model
- Byte buffer reuse through [go-buffer-pool](https://github.com/libp2p/go-buffer-pool)

## Documentation

- [GoDoc](https://godoc.org/github.com/yukino-js/netty-go)

## Repository layout

This repository is a single Go module. Everything is built and tested together with `go build ./...` and `go test ./...`.

| Directory      | Package     | Purpose                                        |
| -------------- | ----------- | ---------------------------------------------- |
| `.`            | `netty`     | Bootstrap, channel, pipeline and handler APIs  |
| `codec/`       | `codec`     | Message codecs                                 |
| `transport/`   | `transport` | Transport interfaces and built-in TCP support  |
| `transport/*/` | protocol    | KCP, QUIC, TLS, UDP and WebSocket transports   |
| `ws/`          | `ws`        | High-level WebSocket server and client API     |
| `internal/`    | internal    | Implementation details, including byte helpers |
| `examples/`    | `main`      | Runnable examples                              |

The core uses `github.com/libp2p/go-buffer-pool` directly for temporary byte slices. Buffer capacities and retention follow libp2p's power-of-two pools.

## Examples

- [chat_server](./examples/chat_server)
- [file_server](./examples/file_server)
- [https_server](./examples/https_server)
- [redis_cli](./examples/redis_cli)
- [tcp_server](./examples/tcp_server)
- [websocket](./examples/websocket)

Run an example from the repository root, for example:

```bash
go run ./examples/tcp_server
```

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
