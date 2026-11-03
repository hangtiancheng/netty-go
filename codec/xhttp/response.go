package xhttp

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

type responseCodec struct {
}

func (*responseCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	reader := utils.MustToReader(message)
	response, err := http.ReadResponse(bufio.NewReader(reader), nil)
	utils.Assert(err)
	ctx.HandleRead(response)
}

func (*responseCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	switch r := message.(type) {
	case *http.Response:
		utils.Assert(r.Write(ctx.Channel().Writer()))
	case ResponseWriter:
	default:
		utils.Assert(fmt.Errorf("unrecognized type: %T", message))
	}
}
