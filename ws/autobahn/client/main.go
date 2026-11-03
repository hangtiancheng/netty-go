package main

import (
	"compress/flate"
	"fmt"
	"log"
	"time"

	netty_ws "github.com/yukino-js/netty-go/ws"
)

const remoteAddr = "127.0.0.1:9001"

func main() {
	const count = 517
	for i := 1; i <= count; i++ {
		testCase(i)
	}
	updateReports()
}

func testCase(id int) {
	var url = fmt.Sprintf("ws://%s/runCase?case=%d&agent=netty_ws/client", remoteAddr, id)
	var onexit = make(chan struct{})

	ws := netty_ws.NewWebsocket(
		netty_ws.WithBufferSize(16*1024, 0),
		netty_ws.WithMaxFrameSize(32*1024*1024),
		netty_ws.WithCompress(flate.BestSpeed, 512),
		netty_ws.WithValidUTF8(),
	)

	ws.OnOpen = func(conn netty_ws.Conn) {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
	}

	ws.OnData = func(conn netty_ws.Conn, data []byte) {
		conn.Write(data)
	}

	ws.OnClose = func(conn netty_ws.Conn, err error) {
		onexit <- struct{}{}
	}

	if _, err := ws.Open(url); nil != err {
		log.Println("ws.Open(", url, ") =>", err)
	}

	<-onexit
}

func updateReports() {
	var url = fmt.Sprintf("ws://%s/updateReports?agent=netty_ws/client", remoteAddr)
	var onexit = make(chan struct{})

	var ws = netty_ws.NewWebsocket(
		netty_ws.WithCompress(flate.BestSpeed, 512),
		netty_ws.WithValidUTF8(),
	)

	ws.OnOpen = func(conn netty_ws.Conn) {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
	}

	ws.OnData = func(conn netty_ws.Conn, data []byte) {
	}

	ws.OnClose = func(conn netty_ws.Conn, err error) {
		onexit <- struct{}{}
	}

	if _, err := ws.Open(url); nil != err {
		log.Println("ws.Open(", url, ") =>", err)
	}

	<-onexit
}
