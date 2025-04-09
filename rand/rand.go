package rand

import (
	"bytes"
	"encoding/binary"
	"github.com/azeroth-sha/simple/buff"
	"io"
	"sync"
)

const DefaultDict = "0123456789abcdefghijklmnopqrstuvwxyz" // 默认字符集

// Random 定义了随机数生成器的接口
type Random interface {
	Int() int
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
	Uint() uint
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Float32() float32
	Float64() float64
	Chars(length int) string
	CharsWith(length int, dict []byte) string
	TextWith(length int, dict []rune) string
}

// random 实现了 Random 接口，用于生成随机数
type random struct {
	pool  *sync.Pool // 用于管理缓冲区的 sync.Pool
	size  int        // 缓冲区大小
	chars []byte     // 默认字符集
}

// Int 生成一个随机的 int 类型整数
func (r *random) Int() int {
	var num int
	if num = int(r.Int64()); num < 0 {
		num = -num
	}
	return num
}

// Int8 生成一个随机的 int8 类型整数
func (r *random) Int8() int8 {
	var num int8
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

// Int16 生成一个随机的 int16 类型整数
func (r *random) Int16() int16 {
	var num int16
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

// Int32 生成一个随机的 int32 类型整数
func (r *random) Int32() int32 {
	var num int32
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

// Int64 生成一个随机的 int64 类型整数
func (r *random) Int64() int64 {
	var num int64
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

// Uint 生成一个随机的 uint 类型整数
func (r *random) Uint() uint {
	var num uint
	num = uint(r.Uint64())
	return num
}

// Uint8 生成一个随机的 uint8 类型整数
func (r *random) Uint8() uint8 {
	var num uint8
	r.read(&num)
	return num
}

// Uint16 生成一个随机的 uint16 类型整数
func (r *random) Uint16() uint16 {
	var num uint16
	r.read(&num)
	return num
}

// Uint32 生成一个随机的 uint32 类型整数
func (r *random) Uint32() uint32 {
	var num uint32
	r.read(&num)
	return num
}

// Uint64 生成一个随机的 uint64 类型整数
func (r *random) Uint64() uint64 {
	var num uint64
	r.read(&num)
	return num
}

// Float32 生成一个随机的 float32 类型浮点数
func (r *random) Float32() float32 {
	var num float32
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

// Float64 生成一个随机的 float64 类型浮点数
func (r *random) Float64() float64 {
	var num float64
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

// Chars 生成指定长度的随机字符串，使用默认字符集
func (r *random) Chars(length int) string {
	return r.CharsWith(length, r.chars)
}

// CharsWith 生成指定长度的随机字符串，使用自定义字符集
func (r *random) CharsWith(length int, dict []byte) string {
	buf := buff.GetBuff()
	defer buff.PutBuff(buf)
	max := uint8(len(dict))
	for i := 0; i < length; i++ {
		_ = buf.WriteByte(dict[r.Uint8()%max])
	}
	return buf.String()
}

// TextWith 生成指定长度的随机文本，使用自定义字符集
func (r *random) TextWith(length int, dict []rune) string {
	buf := buff.GetBuff()
	defer buff.PutBuff(buf)
	max := uint32(len(dict))
	for i := 0; i < length; i++ {
		_, _ = buf.WriteRune(dict[r.Uint32()%max])
	}
	return buf.String()
}

/*
  Package method
*/

// read 从随机数生成器中读取指定类型的随机值
func (r *random) read(v any) {
	buf := r.rndBuff(binary.Size(v))
	defer buff.PutBuff(buf)
	_ = binary.Read(buf, binary.BigEndian, v)
}

// rndBuff 生成指定长度的随机字节缓冲区
func (r *random) rndBuff(n int) *bytes.Buffer {
	buf := buff.GetBuff()
	for rnd := r.getBuff(); buf.Len() < n; buf = r.getBuff() {
		_, _ = io.CopyN(buf, rnd, int64(n-buf.Len()))
		r.putBuff(rnd)
	}
	return buf
}

// putBuff 将缓冲区放回 sync.Pool 中
func (r *random) putBuff(buf *bytes.Buffer) {
	if buf.Len() == 0 {
		buff.PutBuff(buf)
	} else {
		r.pool.Put(buf)
	}
}

// getBuff 从 sync.Pool 中获取一个缓冲区
func (r *random) getBuff() *bytes.Buffer {
	return r.pool.Get().(*bytes.Buffer)
}

// newBuff 创建一个新的随机字节缓冲区
func (r *random) newBuff() interface{} {
	buf := buff.GetBuff()
	if err := readBuff(buf, r.size); err != nil {
		panic(err)
	}
	return buf
}
