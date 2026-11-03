module github.com/yukino-js/netty-go/examples

go 1.26.0

require (
	github.com/yukino-js/netty-go v1.6.8
	github.com/yukino-js/netty-go/transports v1.7.15
)

require github.com/gorilla/websocket v1.5.3 // indirect

replace github.com/yukino-js/netty-go => ../

replace github.com/yukino-js/netty-go/transports => ../transports
