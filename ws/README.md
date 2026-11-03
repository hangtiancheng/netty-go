# netty-go-ws

[![GoDoc][1]][2] [![license-Apache 2][3]][4]

<!--[![Downloads][7]][8]-->

[1]: https://godoc.org/github.com/yukino-js/netty-go/ws?status.svg
[2]: https://godoc.org/github.com/yukino-js/netty-go/ws
[3]: https://img.shields.io/badge/license-Apache%202-blue.svg
[4]: LICENSE

An Websocket server & client written by [netty-go](https://github.com/yukino-js/netty-go)

## Install

```
go get github.com/yukino-js/netty-go/ws@latest
```

## API overview

```
type Websocket
    func NewWebsocket(options ...Option) *Websocket
    func (ws *Websocket) Close() error
    func (ws *Websocket) Listen(addr string) error
    func (ws *Websocket) Open(addr string) (Conn, error)
    func (ws *Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request)
    func (ws *Websocket) UpgradeHTTP(w http.ResponseWriter, r *http.Request) (Conn, error)

type Option
    func WithAsyncWrite(writeQueueSize int, writeForever bool) Option
    func WithBinary() Option
    func WithBufferSize(readBufferSize, writeBufferSize int) Option
    func WithCompress(compressLevel int, compressThreshold int64) Option
    func WithClientHeader(header http.Header) Option
    func WithDialer(dialer Dialer) Option
    func WithMaxFrameSize(maxFrameSize int64) Option
    func WithNoDelay(noDelay bool) Option
    func WithServerHeader(header http.Header) Option
    func WithServeMux(serveMux *http.ServeMux) Option
    func WithServeTLS(tls *tls.Config) Option
    func WithValidUTF8() Option
```

## Easy to use

> Note: `netty_ws` does not support mixed text messages and binary messages, use the `WithBinary` option to switch to binary message mode.

### server :

```go
// create websocket instance
var ws = netty_ws.NewWebsocket()

// setup OnOpen handler
ws.OnOpen = func(conn netty_ws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
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
```

### client :

```go
// create websocket instance
var ws = netty_ws.NewWebsocket()

// setup OnOpen handler
ws.OnOpen = func(conn netty_ws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
    conn.Write([]byte("hello world"))
}

// setup OnData handler
ws.OnData = func(conn netty_ws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
}

// setup OnClose handler
ws.OnClose = func(conn netty_ws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("open websocket connection ...")

// connect to websocket server
if _, err := ws.Open("ws://127.0.0.1:9527/ws"); nil != err {
    panic(err)
}
```

### upgrade from http server:

```go
// create websocket instance
var ws = netty_ws.NewWebsocket()

// setup OnOpen handler
ws.OnOpen = func(conn netty_ws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
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

fmt.Println("upgrade websocket connections ....")

// upgrade websocket connection from http server
http.Handle("/ws", ws)

// listen http server
if err := http.ListenAndServe(":9527", nil); nil != err {
    panic(err)
}
```

## WebSocket backend

The implementation uses `github.com/gorilla/websocket`. The `Conn` interface and `WithClientHeader`, `WithServerHeader`, `WithDialer`, and `WithCompress` options keep their existing signatures. `WithMaxFrameSize` now limits the complete message, including all fragments and decompressed data. Compression uses no context takeover.

HTTP upgrades use Gorilla's default origin policy: absent origins and origins matching the request host are accepted. Applications that need a custom origin policy can use the transport's `HTTPUpgrader` and configure its Gorilla `Upgrader.CheckOrigin`. See the [transport migration notes](../transports/README.md#websocket-migration) for low-level configuration changes.

## Associated

- https://github.com/yukino-js/netty-go
- https://github.com/yukino-js/netty-go/transports
- https://github.com/gorilla/websocket
- https://github.com/lesismal/go-websocket-benchmark
- https://github.com/crossbario/autobahn-testsuite
