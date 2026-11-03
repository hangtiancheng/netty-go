package xhttp

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

type requestCodec struct {
}

func (*requestCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	switch r := message.(type) {
	case io.Reader:
		bufReader := bufio.NewReader(r)
		for {
			request, err := http.ReadRequest(bufReader)
			utils.Assert(err)
			// TODO: replace request context by the channel context
			//
			ctx.HandleRead(request)
			// Close indicates whether to close the connection after
			// replying to this request
			if request.Close {
				ctx.Close(fmt.Errorf("request mark closed"))
				return
			}
		}
	default:
		ctx.HandleRead(message)
	}
}

func (*requestCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	switch r := message.(type) {
	case *http.Request:
		utils.Assert(r.Write(ctx.Channel().Writer()))
	default:
		ctx.HandleWrite(message)
	}
}
