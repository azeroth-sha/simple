package codec

type Type uint16

const (
	Json Type = iota + 1 // json
	MsgP                 // msgpack
)
