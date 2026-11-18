package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hangtiancheng/netty-go/ws"
)

func main() {

	header := http.Header{}
	header.Add("client-time", time.Now().String())

	// create websocket instance
	client := ws.NewWebsocket(ws.WithClientHeader(header))

	// setup OnOpen handler
	client.OnOpen = func(conn ws.Conn) {
		fmt.Println("OnOpen: ", conn.RemoteAddr(), ", header: ", conn.Header())
		conn.Write([]byte("hello world"))
	}

	// setup OnData handler
	client.OnData = func(conn ws.Conn, data []byte) {
		fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
	}

	// setup OnClose handler
	client.OnClose = func(conn ws.Conn, err error) {
		fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
	}

	fmt.Println("open websocket connection ...")

	// connect to websocket server
	if _, err := client.Open("ws://127.0.0.1:9527/ws"); nil != err {
		panic(err)
	}

	select {}
}
