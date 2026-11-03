package redis_go

import "bytes"

type RespKind byte

const (
	SimpleKind  RespKind = '+'
	ErrorKind   RespKind = '-'
	IntegerKind RespKind = ':'
	BlukKind    RespKind = '$'
	ArrayKind   RespKind = '*'
)

type Resp struct {
	Kind  RespKind
	Null  bool
	Data  string
	Array []Resp
}

func (r *Resp) String() string {
	w := &bytes.Buffer{}
	EncodeResp(w, r)
	return w.String()
}
