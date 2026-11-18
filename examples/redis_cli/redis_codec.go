package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hangtiancheng/netty-go"
	"github.com/hangtiancheng/netty-go/examples/redis_cli/resp"
)

type simpleRedisCodec struct {
	decoder *resp.Decoder
}

func (s *simpleRedisCodec) CodecName() string {
	return "simple-redis-codec"
}

func (s *simpleRedisCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	// init decoder.
	if nil == s.decoder {
		s.decoder = resp.NewDecoder(message.(io.Reader), 10240)
	}

	// decode redis response.
	reply := &resp.Resp{}
	if err := s.decoder.Decode(reply); nil != err {
		ctx.Close(err)
		return
	}

	// fmt.Println("->", reply.String())
	// post response.
	ctx.HandleRead(reply)
}

func (s *simpleRedisCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {

	switch v := message.(type) {
	case []resp.Value:
		// encode request.
		buffer := bytes.NewBuffer(nil)
		if err := resp.EncodeMulti(buffer, v...); nil != err {
			ctx.Close(err)
			return
		}

		// fmt.Println("<-", buffer.String())
		// post request.
		ctx.HandleWrite(buffer)
	default:
		ctx.Close(fmt.Errorf("%T is invalid message", message))
	}
}
