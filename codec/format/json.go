package format

import (
	"encoding/json"

	"github.com/yukino-js/netty-go"
	"github.com/yukino-js/netty-go/codec"
	"github.com/yukino-js/netty-go/utils"
)

// JSONCodec create a json codec
func JSONCodec(useNumber bool, disAllowUnknownFields bool) codec.Codec {
	return &jsonCodec{
		useNumber:             useNumber,
		disAllowUnknownFields: disAllowUnknownFields,
	}
}

type jsonCodec struct {
	useNumber             bool
	disAllowUnknownFields bool
}

func (*jsonCodec) CodecName() string {
	return "json-codec"
}

func (j *jsonCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {

	// new json decoder from reader
	jsonDecoder := json.NewDecoder(utils.MustToReader(message))

	// UseNumber causes the Decoder to unmarshal a number into an any as a
	// Number instead of as a float64.
	if j.useNumber {
		jsonDecoder.UseNumber()
	}

	// DisallowUnknownFields causes the Decoder to return an error when the destination
	// is a struct and the input contains object keys which do not match any
	// non-ignored, exported fields in the destination.
	if j.disAllowUnknownFields {
		jsonDecoder.DisallowUnknownFields()
	}

	// decode to map
	var object = make(map[string]any)
	utils.Assert(jsonDecoder.Decode(&object))

	// post object
	ctx.HandleRead(object)
}

func (j *jsonCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	// marshal object to json bytes
	data := utils.AssertBytes(json.Marshal(message))
	// post json
	ctx.HandleWrite(data)
}
