package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/examples/redis_cli/redis_go"
	"github.com/yukino-js/netty-go/utils"
)

type simpleRedisCodec struct {
	decoder *redis_go.Decoder
}

func (s *simpleRedisCodec) CodecName() string {
	return "simple-redis-codec"
}

func (s *simpleRedisCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	// init decoder.
	if nil == s.decoder {
		s.decoder = redis_go.NewDecoder(message.(io.Reader), 10240)
	}

	// decode redis response.
	resp := &redis_go.Resp{}
	utils.Assert(s.decoder.Decode(resp))

	// fmt.Println("->", resp.String())
	// post response.
	ctx.HandleRead(resp)
}

func (s *simpleRedisCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	switch v := message.(type) {
	case []redis_go.Value:
		// encode request.
		buffer := bytes.NewBuffer(nil)
		utils.Assert(redis_go.EncodeMulti(buffer, v...))

		// fmt.Println("<-", buffer.String())
		// post request.
		ctx.HandleWrite(buffer)
	default:
		utils.Assert(fmt.Errorf("%T is invalid message", message))
	}
}
