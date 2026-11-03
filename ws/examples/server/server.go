package main

import (
	"fmt"
	"net/http"
	"time"

	netty_ws "github.com/yukino-js/netty-go/ws"
)

func main() {
	header := http.Header{}
	header.Add("server-time", time.Now().String())

	// create websocket instance
	var ws = netty_ws.NewWebsocket(netty_ws.WithServerHeader(header))

	// setup OnOpen handler
	ws.OnOpen = func(conn netty_ws.Conn) {
		fmt.Println("OnOpen: ", conn.RemoteAddr(), ", header: ", conn.Header())
	}

	// setup OnData handler
	ws.OnData = func(conn netty_ws.Conn, data []byte) {
		fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
		conn.Write(data)
	}

	// setup OnClose handler
	ws.OnClose = func(conn netty_ws.Conn, err error) {
		fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
	}

	fmt.Println("listening websocket connections ....")
	// listen websocket server
	if err := ws.Listen("ws://127.0.0.1:9527/ws"); nil != err {
		panic(err)
	}
}
