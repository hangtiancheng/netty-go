package main

import (
	"compress/flate"
	"fmt"
	"log"
	"time"

	"github.com/hangtiancheng/netty-go/ws"
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
	var url = fmt.Sprintf("ws://%s/runCase?case=%d&agent=netty-go/client", remoteAddr, id)
	var onexit = make(chan struct{})

	client := ws.NewWebsocket(
		ws.WithBufferSize(16*1024, 0),
		ws.WithMaxFrameSize(32*1024*1024),
		ws.WithCompress(flate.BestSpeed, 512),
		ws.WithValidUTF8(),
	)

	client.OnOpen = func(conn ws.Conn) {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
	}

	client.OnData = func(conn ws.Conn, data []byte) {
		conn.Write(data)
	}

	client.OnClose = func(conn ws.Conn, err error) {
		onexit <- struct{}{}
	}

	if _, err := client.Open(url); nil != err {
		log.Println("client.Open(", url, ") =>", err)
	}

	<-onexit
}

func updateReports() {
	var url = fmt.Sprintf("ws://%s/updateReports?agent=netty-go/client", remoteAddr)
	var onexit = make(chan struct{})

	client := ws.NewWebsocket(
		ws.WithCompress(flate.BestSpeed, 512),
		ws.WithValidUTF8(),
	)

	client.OnOpen = func(conn ws.Conn) {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
	}

	client.OnData = func(conn ws.Conn, data []byte) {
	}

	client.OnClose = func(conn ws.Conn, err error) {
		onexit <- struct{}{}
	}

	if _, err := client.Open(url); nil != err {
		log.Println("client.Open(", url, ") =>", err)
	}

	<-onexit
}
