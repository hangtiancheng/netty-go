package netty_ws

import (
	"strconv"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/examples/transport/websocket"
)

// ClosedError returned when peer has closed the connection with appropriate
// code and a textual reason.
type ClosedError struct {
	Code   int
	Reason string
}

// Error implements error interface.
func (err ClosedError) Error() string {
	return "ws closed: " + strconv.FormatUint(uint64(err.Code), 10) + " " + err.Reason
}

// ErrServerClosed is returned by the Server call Shutdown or Close
var ErrServerClosed = netty.ErrServerClosed

var defaultEngine = netty.NewBootstrap(
	netty.WithTransport(websocket.New()),
	netty.WithChannel(netty.NewChannel()),
	netty.WithChannelHolder(nil),
	netty.WithClientInitializer(makeInitializer(true)),
	netty.WithChildInitializer(makeInitializer(false)),
)

func makeInitializer(client bool) netty.ChannelInitializer {
	return func(channel netty.Channel) {
		ws := channel.Attachment().(*Websocket)
		channel.Pipeline().
			AddLast(ws.holder).
			AddLast(newConn(ws, channel, client))
	}
}
