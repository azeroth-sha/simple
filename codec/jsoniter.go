//go:build jsoniter
// +build jsoniter

package codec

import (
	jsoniter "github.com/json-iterator/go"
)

func JsonMarshal(v any) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func JsonMustMarshal(v any) []byte {
	b, e := JsonMarshal(v)
	if e != nil {
		panic(e)
	}
	return b
}

func JsonUnmarshal(b []byte, v any) error {
	return jsoniter.Unmarshal(b, v)
}

func JsonMustUnmarshal(b []byte, v any) {
	if e := JsonUnmarshal(b, v); e != nil {
		panic(e)
	}
}
