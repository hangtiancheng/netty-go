# netty-go transports

[![GoDoc][1]][2] [![license-Apache 2][3]][4]

[1]: https://godoc.org/github.com/hangtiancheng/netty-go/transport?status.svg
[2]: https://godoc.org/github.com/hangtiancheng/netty-go/transport
[3]: https://img.shields.io/badge/license-Apache%202-blue.svg
[4]: ../LICENSE

Transport abstractions and implementations for [netty-go](https://github.com/hangtiancheng/netty-go).

The `transport` package defines the `Transport`, `Acceptor` and `Factory` interfaces. All implementations are part of the same Go module and live in subdirectories:

| Directory             | Scheme                | Notes                                      |
| --------------------- | --------------------- | ------------------------------------------ |
| `transport/tcp`       | `tcp`, `tcp4`, `tcp6` | Built-in default transport                 |
| `transport/udp`       | `udp`, `udp4`, `udp6` | Datagram transport with optional reuseport |
| `transport/kcp`       | `kcp`                 | KCP reliable UDP transport                 |
| `transport/quic`      | `quic`                | QUIC transport                             |
| `transport/tls`       | `tcp`, `tcp4`, `tcp6` | TLS over TCP                               |
| `transport/websocket` | `ws`, `wss`           | WebSocket transport                        |

The high-level WebSocket server and client API is available in [`../ws`](../ws).

## WebSocket transport

The WebSocket transport uses `github.com/gorilla/websocket`. Handshake parsing, masking, fragmentation, control frames, and per-message compression are handled by Gorilla.

`Options.Apply()` returns a configured copy, so callers should retain its return value when calling it directly. `WithOptions` applies the configuration automatically.

Compression retains the configured level and write threshold, and only runs when negotiated. Gorilla supports `permessage-deflate` with no context takeover. `ReadBufferSize` and `WriteBufferSize` configure Gorilla's buffers. Each write completes a message immediately, so `Flush()` is a no-op. Reads return `io.EOF` at each message boundary.

`MaxFrameSize` limits an entire incoming message across all fragments, including its decompressed payload. Zero means no limit. `CheckUTF8` validates text messages across reads and fragments.

Gorilla's default origin policy allows requests with no `Origin`, or an `Origin` matching the request host. Configure `Options.Upgrader.CheckOrigin` explicitly when other origins must be allowed. See the [Gorilla API documentation](https://pkg.go.dev/github.com/gorilla/websocket).
