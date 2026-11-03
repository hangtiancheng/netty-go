module github.com/yukino-js/netty-go/transports

go 1.26.0

require (
	github.com/gorilla/websocket v1.5.3
	github.com/libp2p/go-reuseport v0.4.0
	github.com/quic-go/quic-go v0.65.0
	github.com/xtaci/kcp-go/v5 v5.7.2
	github.com/yukino-js/netty-go v1.6.8
	golang.org/x/crypto v0.58.0
)

require (
	github.com/klauspost/cpuid/v2 v2.4.2 // indirect
	github.com/klauspost/reedsolomon v1.14.2 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/time v0.16.0 // indirect
)

replace github.com/yukino-js/netty-go => ../
