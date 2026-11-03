# netty-go-transports

[![GoDoc][1]][2] [![license-Apache 2][3]][4]

[1]: https://godoc.org/github.com/yukino-js/netty-go/transports?status.svg
[2]: https://godoc.org/github.com/yukino-js/netty-go/transports
[3]: https://img.shields.io/badge/license-Apache%202-blue.svg
[4]: LICENSE

Open source transport implementations for [netty-go](https://github.com/yukino-js/netty-go). This is a standalone module in the netty-go workspace.

## Transport

1. TCP (builtin)
2. KCP (A Fast and Reliable ARQ Protocol)
3. Websocket
4. QUIC (Quick UDP Internet Connections)
5. UDP
6. TLS (Transport Layer Security over Tcp)

## WebSocket migration

The WebSocket transport uses `github.com/gorilla/websocket` v1.5.3. Handshake parsing, masking, fragmentation, control frames, and per-message compression are handled by Gorilla.

Custom transport configurations need these changes:

| Previous configuration                 | Current configuration                                              |
| -------------------------------------- | ------------------------------------------------------------------ |
| `ws.OpText`, `ws.OpBinary`             | `websocket.TextMessage`, `websocket.BinaryMessage` from Gorilla    |
| `Options.Dialer` (`ws.Dialer`)         | `websocket.Dialer`; use `HandshakeTimeout` and `NetDialContext`    |
| `Options.Upgrader` (`ws.HTTPUpgrader`) | `websocket.Upgrader`                                               |
| `Dialer.Header` / `Upgrader.Header`    | `Options.RequestHeader` / `Options.ResponseHeader` (`http.Header`) |
| Custom extension negotiation           | `Options.CompressEnabled`                                          |

`Options.Apply()` returns a configured copy; retain its return value when calling it directly. `WithOptions` applies the configuration automatically.

Compression retains the configured level and write threshold, and only runs when negotiated. Gorilla supports `permessage-deflate` with no context takeover; custom window sizes and context takeover are no longer exposed. `ReadBufferSize` and `WriteBufferSize` configure Gorilla's buffers. Each write completes a message immediately, so `Flush()` is a no-op. Reads still return `io.EOF` at each message boundary.

`MaxFrameSize` now limits an entire incoming message across all fragments, including its decompressed payload. Zero means no limit. `CheckUTF8` still validates text messages across reads and fragments.

Gorilla's default origin policy allows requests with no `Origin`, or an `Origin` matching the request host. Configure `Options.Upgrader.CheckOrigin` explicitly when other origins must be allowed. See the [Gorilla API documentation](https://pkg.go.dev/github.com/gorilla/websocket).
