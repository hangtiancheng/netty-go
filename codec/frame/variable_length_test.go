package frame

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

func TestVariableLengthCodec(t *testing.T) {

	var cases = []struct {
		maxReadLen int
		input      []byte
		output     any
	}{
		{maxReadLen: 1024, output: "123456789"},
		{maxReadLen: 1024, output: strings.NewReader("123456789")},
	}

	for index, c := range cases {
		codec := VariableLengthCodec(c.maxReadLen)
		t.Run(fmt.Sprint(codec.CodecName(), "#", index), func(t *testing.T) {
			ctx := MockHandlerContext{
				MockHandleRead: func(message netty.Message) {
					if dst := utils.MustToBytes(message); !bytes.Equal(dst, c.input) {
						t.Fatalf("%v != %v", dst, c.input)
					}
				},

				MockHandleWrite: func(message netty.Message) {
					c.input = utils.MustToBytes(message)
				},
			}
			codec.HandleWrite(ctx, c.output)
			codec.HandleRead(ctx, c.input)
		})
	}

}
