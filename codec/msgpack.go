package codec

import (
	"github.com/vmihailenco/msgpack/v5"
)

func MsgPMarshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func MsgPMustMarshal(v any) []byte {
	b, e := JsonMarshal(v)
	if e != nil {
		panic(e)
	}
	return b
}

func MsgPUnmarshal(b []byte, v any) error {
	return msgpack.Unmarshal(b, v)
}

func MsgPMustUnmarshal(b []byte, v any) {
	if e := JsonUnmarshal(b, v); e != nil {
		panic(e)
	}
}
