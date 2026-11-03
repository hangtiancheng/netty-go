package frame

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/yukino-js/netty-go/utils"

	"github.com/yukino-js/netty-go"
)

func TestLengthFieldCodec(t *testing.T) {

	var cases = []struct {
		byteOrder        binary.ByteOrder
		maxFrameLen      int
		fieldOffset      int
		fieldLen         int
		lengthAdjustment int
		bytesToStrip     int
		input            []byte
		output           any
	}{
		{byteOrder: binary.LittleEndian, maxFrameLen: 1024, fieldLen: 1, lengthAdjustment: 0, bytesToStrip: 1, output: []byte("123456789")},
		{byteOrder: binary.LittleEndian, maxFrameLen: 1024, fieldLen: 2, lengthAdjustment: 0, bytesToStrip: 2, output: []byte("123456789")},
		{byteOrder: binary.LittleEndian, maxFrameLen: 1024, fieldLen: 4, lengthAdjustment: 0, bytesToStrip: 4, output: []byte("123456789")},
		{byteOrder: binary.LittleEndian, maxFrameLen: 1024, fieldLen: 8, lengthAdjustment: 0, bytesToStrip: 8, output: []byte("123456789")},

		{byteOrder: binary.BigEndian, maxFrameLen: 1024, fieldLen: 1, lengthAdjustment: 0, bytesToStrip: 1, output: []byte("123456789")},
		{byteOrder: binary.BigEndian, maxFrameLen: 1024, fieldLen: 2, lengthAdjustment: 0, bytesToStrip: 2, output: []byte("123456789")},
		{byteOrder: binary.BigEndian, maxFrameLen: 1024, fieldLen: 4, lengthAdjustment: 0, bytesToStrip: 4, output: []byte("123456789")},
		{byteOrder: binary.BigEndian, maxFrameLen: 1024, fieldLen: 8, lengthAdjustment: 0, bytesToStrip: 8, output: []byte("123456789")},
	}

	for _, c := range cases {
		codec := LengthFieldCodec(c.byteOrder, c.maxFrameLen, c.fieldOffset, c.fieldLen, c.lengthAdjustment, c.bytesToStrip)
		t.Run(codec.CodecName(), func(t *testing.T) {
			ctx := MockHandlerContext{
				MockHandleRead: func(message netty.Message) {
					if dst := utils.MustToBytes(message); !bytes.Equal(dst, c.input[c.bytesToStrip:]) {
						t.Fatalf("%v != %v", dst, c.input[c.bytesToStrip:])
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
