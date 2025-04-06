package codec

import (
	code "github.com/ugorji/go/codec"
)

// Marshal encode to bytes
func Marshal(t Type, v any) ([]byte, error) {
	b := make([]byte, 0)
	e := code.NewEncoderBytes(&b, getHandler(t)).
		Encode(v)
	return b, e
}

// MustMarshal encode to bytes
func MustMarshal(t Type, v any) []byte {
	b := make([]byte, 0)
	code.NewEncoderBytes(&b, getHandler(t)).
		MustEncode(v)
	return b
}

// Unmarshal decode from bytes
func Unmarshal(t Type, b []byte, v any) error {
	return code.NewDecoderBytes(b, getHandler(t)).
		Decode(v)
}

// MustUnmarshal decode from bytes
func MustUnmarshal(t Type, b []byte, v any) {
	code.NewDecoderBytes(b, getHandler(t)).
		MustDecode(v)
}
