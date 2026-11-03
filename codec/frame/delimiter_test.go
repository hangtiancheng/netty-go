package frame

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yukino-js/netty-go/utils"

	"github.com/yukino-js/netty-go"
)

func TestDelimiterCodec(t *testing.T) {

	var cases = []struct {
		delimiter      string
		stripDelimiter bool
		input          []byte
		output         any
	}{
		{delimiter: "\n", stripDelimiter: true, output: []byte("123456789")},
		{delimiter: "\n", stripDelimiter: false, output: []byte("123456789")},
		{delimiter: "\n", stripDelimiter: true, output: []byte("123456789")},
		{delimiter: "\n", stripDelimiter: false, output: []byte("123456789")},
		{delimiter: "\r\n", stripDelimiter: true, output: []byte("123456789")},
		{delimiter: "\r\n", stripDelimiter: false, output: []byte("123456789")},
		{delimiter: "\r\n", stripDelimiter: true, output: []byte("123456789")},
		{delimiter: "\r\n", stripDelimiter: false, output: []byte("123456789")},
		{delimiter: "\n", stripDelimiter: false, output: strings.NewReader("123456789")},
		{delimiter: "\r\n", stripDelimiter: true, output: strings.NewReader("123456789")},
	}

	for index, c := range cases {
		codec := DelimiterCodec(1024, c.delimiter, c.stripDelimiter)
		t.Run(fmt.Sprint(codec.CodecName(), "#", index), func(t *testing.T) {
			ctx := MockHandlerContext{
				MockHandleRead: func(message netty.Message) {
					if c.stripDelimiter {
						c.input = bytes.TrimRight(c.input, c.delimiter)
					}
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
