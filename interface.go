package simple

// String 定义了返回字符串表示的接口
type String interface {
	String() string // 返回对象的字符串表示
}

// Bytes 定义了返回字节切片的接口
type Bytes interface {
	Bytes() []byte // 返回对象的字节切片表示
}

// Int 定义了返回 int 类型值的接口
type Int interface {
	Int() int // 返回对象的 int 类型值
}

// Int8 定义了返回 int8 类型值的接口
type Int8 interface {
	Int8() int8 // 返回对象的 int8 类型值
}

// Int16 定义了返回 int16 类型值的接口
type Int16 interface {
	Int16() int16 // 返回对象的 int16 类型值
}

// Int32 定义了返回 int32 类型值的接口
type Int32 interface {
	Int32() int32 // 返回对象的 int32 类型值
}

// Int64 定义了返回 int64 类型值的接口
type Int64 interface {
	Int64() int64 // 返回对象的 int64 类型值
}

// Uint 定义了返回 uint 类型值的接口
type Uint interface {
	Uint() uint // 返回对象的 uint 类型值
}

// Uint8 定义了返回 uint8 类型值的接口
type Uint8 interface {
	Uint8() uint8 // 返回对象的 uint8 类型值
}

// Uint16 定义了返回 uint16 类型值的接口
type Uint16 interface {
	Uint16() uint16 // 返回对象的 uint16 类型值
}

// Uint32 定义了返回 uint32 类型值的接口
type Uint32 interface {
	Uint32() uint32 // 返回对象的 uint32 类型值
}

// Uint64 定义了返回 uint64 类型值的接口
type Uint64 interface {
	Uint64() uint64 // 返回对象的 uint64 类型值
}

// Float32 定义了返回 float32 类型值的接口
type Float32 interface {
	Float32() float32 // 返回对象的 float32 类型值
}

// Float64 定义了返回 float64 类型值的接口
type Float64 interface {
	Float64() float64 // 返回对象的 float64 类型值
}

// Bool 定义了返回 bool 类型值的接口
type Bool interface {
	Bool() bool // 返回对象的 bool 类型值
}
