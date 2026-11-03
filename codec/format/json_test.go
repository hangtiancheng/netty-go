package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/utils"
)

func TestJsonCodec_HandleWrite(t *testing.T) {

	var cases = []struct {
		useNumber             bool
		disAllowUnknownFields bool
		input                 []byte
		output                any
	}{
		{useNumber: true, disAllowUnknownFields: false, output: map[string]any{"key": 123456789}},
		{useNumber: false, disAllowUnknownFields: false, output: map[string]any{"key": 123456789}},
		{useNumber: false, disAllowUnknownFields: true, output: map[string]any{"key": 123456789}},
	}

	for index, c := range cases {
		codec := JSONCodec(c.useNumber, c.disAllowUnknownFields)
		t.Run(fmt.Sprint(codec.CodecName(), "#", index), func(t *testing.T) {
			ctx := MockHandlerContext{
				MockHandleRead: func(message netty.Message) {
					if dst := utils.AssertBytes(json.Marshal(message)); !bytes.Equal(dst, c.input) {
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
