package main

import (
	"fmt"
	"net/http"

	"github.com/yukino-js/netty-go/ws"
)

func main() {

	// create websocket instance
	server := ws.NewWebsocket()

	// setup OnOpen handler
	server.OnOpen = func(conn ws.Conn) {
		fmt.Println("OnOpen: ", conn.RemoteAddr())
	}

	// setup OnData handler
	server.OnData = func(conn ws.Conn, data []byte) {
		fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data), ", header: ", conn.Header())
		conn.Write(data)
	}

	// setup OnClose handler
	server.OnClose = func(conn ws.Conn, err error) {
		fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
	}

	fmt.Println("upgrade websocket connections ....")

	// upgrade websocket connection from http server
	http.Handle("/ws", server)

	// listen http server
	if err := http.ListenAndServe(":9527", nil); nil != err {
		panic(err)
	}
}
