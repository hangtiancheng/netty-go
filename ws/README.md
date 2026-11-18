# netty-go ws

[![GoDoc][1]][2] [![license-Apache 2][3]][4]

[1]: https://godoc.org/github.com/yukino-js/netty-go/ws?status.svg
[2]: https://godoc.org/github.com/yukino-js/netty-go/ws
[3]: https://img.shields.io/badge/license-Apache%202-blue.svg
[4]: ../LICENSE

A high-level WebSocket server and client built on [netty-go](https://github.com/yukino-js/netty-go).

The low-level WebSocket transport lives in [`../transport/websocket`](../transport/websocket).

## Install

```
go get github.com/yukino-js/netty-go
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
    func WithDialTimeout(timeout time.Duration) Option
    func WithMaxFrameSize(maxFrameSize int64) Option
    func WithNoDelay(noDelay bool) Option
    func WithServerHeader(header http.Header) Option
    func WithServeMux(serveMux *http.ServeMux) Option
    func WithServeTLS(tls *tls.Config) Option
    func WithValidUTF8() Option
```

## Easy to use

> Note: `ws` does not support mixed text and binary messages in the same connection. Use the `WithBinary` option to switch to binary message mode.

### server

```go
server := ws.NewWebsocket()

server.OnOpen = func(conn ws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
}

server.OnData = func(conn ws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
    conn.Write(data)
}

server.OnClose = func(conn ws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("listening websocket connections ....")

if err := server.Listen("ws://127.0.0.1:9527/ws"); nil != err {
    panic(err)
}
```

### client

```go
client := ws.NewWebsocket()

client.OnOpen = func(conn ws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
    conn.Write([]byte("hello world"))
}

client.OnData = func(conn ws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
}

client.OnClose = func(conn ws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("open websocket connection ...")

if _, err := client.Open("ws://127.0.0.1:9527/ws"); nil != err {
    panic(err)
}
```

### upgrade from an HTTP server

```go
server := ws.NewWebsocket()

server.OnOpen = func(conn ws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
}

server.OnData = func(conn ws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
    conn.Write(data)
}

server.OnClose = func(conn ws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("upgrade websocket connections ....")

http.Handle("/ws", server)

if err := http.ListenAndServe(":9527", nil); nil != err {
    panic(err)
}
```

## WebSocket backend

The implementation uses `github.com/gorilla/websocket`. `WithMaxFrameSize` limits the complete message, including all fragments and decompressed data. Compression uses no context takeover.

HTTP upgrades use Gorilla's default origin policy: absent origins and origins matching the request host are accepted. Applications that need a custom origin policy can use the transport's `HTTPUpgrader` and configure its Gorilla `Upgrader.CheckOrigin`. See the [transport documentation](../transport/README.md#websocket-transport) for low-level configuration details.

## Associated

- https://github.com/yukino-js/netty-go
- https://github.com/yukino-js/netty-go/tree/main/transport/websocket
- https://github.com/gorilla/websocket
- https://github.com/lesismal/go-websocket-benchmark
- https://github.com/crossbario/autobahn-testsuite
