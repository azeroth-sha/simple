package codec

import (
	"errors"
)

var ErrUnknownType = errors.New(`unknown type`)

type Type uint16

const (
	Json Type = iota + 1 // json
	MsgP                 // msgpack
)

// Marshal encode to bytes
func Marshal(t Type, v any) ([]byte, error) {
	switch t {
	case Json:
		return JsonMarshal(v)
	case MsgP:
		return MsgPMarshal(v)
	default:
		return nil, ErrUnknownType
	}
}

// MustMarshal encode to bytes
func MustMarshal(t Type, v any) []byte {
	switch t {
	case Json:
		return JsonMustMarshal(v)
	case MsgP:
		return MsgPMustMarshal(v)
	default:
		panic(ErrUnknownType)
	}
}

// Unmarshal decode from bytes
func Unmarshal(t Type, b []byte, v any) error {
	switch t {
	case Json:
		return JsonUnmarshal(b, v)
	case MsgP:
		return MsgPUnmarshal(b, v)
	default:
		return ErrUnknownType
	}
}

// MustUnmarshal decode from bytes
func MustUnmarshal(t Type, b []byte, v any) {
	switch t {
	case Json:
		JsonMustUnmarshal(b, v)
	case MsgP:
		MsgPMustUnmarshal(b, v)
	default:
		panic(ErrUnknownType)
	}
}
