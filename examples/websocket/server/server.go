package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yukino-js/netty-go/ws"
)

func main() {
	header := http.Header{}
	header.Add("server-time", time.Now().String())

	// create websocket instance
	server := ws.NewWebsocket(ws.WithServerHeader(header))

	// setup OnOpen handler
	server.OnOpen = func(conn ws.Conn) {
		fmt.Println("OnOpen: ", conn.RemoteAddr(), ", header: ", conn.Header())
	}

	// setup OnData handler
	server.OnData = func(conn ws.Conn, data []byte) {
		fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
		conn.Write(data)
	}

	// setup OnClose handler
	server.OnClose = func(conn ws.Conn, err error) {
		fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
	}

	fmt.Println("listening websocket connections ....")
	// listen websocket server
	if err := server.Listen("ws://127.0.0.1:9527/ws"); nil != err {
		panic(err)
	}
}
