package main

import (
	"github.com/yukino-js/netty-go/ws"
)

func main() {

	// tcpkali -c 1000 --connect-rate 500 -r 1000 -T 30s -f 1K.txt --ws 127.0.0.1:8000
	server := ws.NewWebsocket()

	server.OnData = func(conn ws.Conn, data []byte) {
		_ = conn.Write(data)
	}

	if err := server.Listen(":8000"); nil != err {
		panic(err)
	}
}
