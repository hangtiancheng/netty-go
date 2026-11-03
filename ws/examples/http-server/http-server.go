package main

import (
	"fmt"
	"net/http"

	netty_ws "github.com/yukino-js/netty-go/ws"
)

func main() {

	// create websocket instance
	var ws = netty_ws.NewWebsocket()

	// setup OnOpen handler
	ws.OnOpen = func(conn netty_ws.Conn) {
		fmt.Println("OnOpen: ", conn.RemoteAddr())
	}

	// setup OnData handler
	ws.OnData = func(conn netty_ws.Conn, data []byte) {
		fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data), ", header: ", conn.Header())
		conn.Write(data)
	}

	// setup OnClose handler
	ws.OnClose = func(conn netty_ws.Conn, err error) {
		fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
	}

	fmt.Println("upgrade websocket connections ....")

	// upgrade websocket connection from http server
	http.Handle("/ws", ws)

	// listen http server
	if err := http.ListenAndServe(":9527", nil); nil != err {
		panic(err)
	}
}
