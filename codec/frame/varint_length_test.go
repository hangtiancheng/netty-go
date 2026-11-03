package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

func TestVarintLengthFieldCodec(t *testing.T) {

	var cases = []struct {
		maxFrameLen int
		input       []byte
		output      any
		offset      int
	}{
		{maxFrameLen: 1024, output: "123456789"},
		{maxFrameLen: 1024, output: strings.NewReader("123456789")},
	}

	for index, c := range cases {
		codec := VarintLengthFieldCodec(c.maxFrameLen)
		t.Run(fmt.Sprint(codec.CodecName(), "#", index), func(t *testing.T) {
			ctx := MockHandlerContext{
				MockHandleRead: func(message netty.Message) {
					if dst := utils.MustToBytes(message); !bytes.Equal(dst, c.input[c.offset:]) {
						t.Fatalf("%v != %v", dst, c.input[c.offset:])
					}
				},

				MockHandleWrite: func(message netty.Message) {
					c.input = utils.MustToBytes(message)
					_, c.offset = binary.Uvarint(c.input)
				},
			}
			codec.HandleWrite(ctx, c.output)
			codec.HandleRead(ctx, c.input)
		})
	}
}
