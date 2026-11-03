package main

import (
	netty_ws "github.com/yukino-js/netty-go/ws"
)

func main() {

	// tcpkali -c 1000 --connect-rate 500 -r 1000 -T 30s -f 1K.txt --ws 127.0.0.1:8000
	ws := netty_ws.NewWebsocket()

	ws.OnData = func(conn netty_ws.Conn, data []byte) {
		_ = conn.Write(data)
	}

	if err := ws.Listen(":8000"); nil != err {
		panic(err)
	}
}
