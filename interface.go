package simple

type String interface {
	String() string
}

type Bytes interface {
	Bytes() []byte
}

type Int interface {
	Int() int
}

type Int8 interface {
	Int8() int8
}

type Int16 interface {
	Int16() int16
}

type Int32 interface {
	Int32() int32
}

type Int64 interface {
	Int64() int64
}

type Uint interface {
	Uint() uint
}

type Uint8 interface {
	Uint8() uint8
}

type Uint16 interface {
	Uint16() uint16
}

type Uint32 interface {
	Uint32() uint32
}

type Uint64 interface {
	Uint64() uint64
}

type Float32 interface {
	Float32() float32
}

type Float64 interface {
	Float64() float64
}

type Bool interface {
	Bool() bool
}
