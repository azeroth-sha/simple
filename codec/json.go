package codec

import (
	"github.com/bytedance/sonic"
)

func JsonMarshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

func JsonMustMarshal(v any) []byte {
	b, e := JsonMarshal(v)
	if e != nil {
		panic(e)
	}
	return b
}

func JsonUnmarshal(b []byte, v any) error {
	return sonic.Unmarshal(b, v)
}

func JsonMustUnmarshal(b []byte, v any) {
	if e := JsonUnmarshal(b, v); e != nil {
		panic(e)
	}
}
