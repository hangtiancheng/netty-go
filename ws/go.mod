module github.com/yukino-js/netty-go/ws

go 1.26.0

require (
	github.com/gobwas/ws v1.4.0
	github.com/yukino-js/netty-go v1.6.8
	github.com/yukino-js/netty-go/transports v1.7.15
)

require (
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/klauspost/compress v1.20.3 // indirect
	golang.org/x/sys v0.49.0 // indirect
)

replace github.com/yukino-js/netty-go => ../

replace github.com/yukino-js/netty-go/transports => ../transports
