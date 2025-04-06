package codec

import code "github.com/ugorji/go/codec"

func getHandler(t Type) code.Handle {
	switch t {
	case Json:
		return new(code.JsonHandle)
	case MsgP:
		return new(code.MsgpackHandle)
	default:
		return nil
	}
}
