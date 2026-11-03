package frame

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yukino-js/netty-go/utils"

	"github.com/yukino-js/netty-go"
)

func TestFixedLengthCodec(t *testing.T) {

	var cases = []struct {
		length int
		input  []byte
		output any
	}{
		{length: 5, output: []byte("12345")},
		{length: 5, output: strings.NewReader("12345")},
	}

	for index, c := range cases {
		codec := FixedLengthCodec(c.length)
		t.Run(fmt.Sprint(codec.CodecName(), "#", index), func(t *testing.T) {
			ctx := MockHandlerContext{
				MockHandleRead: func(message netty.Message) {
					if dst := utils.MustToBytes(message); !bytes.Equal(dst, c.input) {
						t.Fatal(dst, "!=", c.input)
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
