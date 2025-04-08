package rand

import (
	"bytes"
	"encoding/binary"
	"github.com/azeroth-sha/simple/buff"
	"io"
	"sync"
)

const DefaultDict = "0123456789abcdefghijklmnopqrstuvwxyz"

// Random interface
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

type random struct {
	pool  *sync.Pool
	size  int
	chars []byte
}

func (r *random) Int() int {
	var num int
	if num = int(r.Int64()); num < 0 {
		num = -num
	}
	return num
}

func (r *random) Int8() int8 {
	var num int8
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

func (r *random) Int16() int16 {
	var num int16
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

func (r *random) Int32() int32 {
	var num int32
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

func (r *random) Int64() int64 {
	var num int64
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

func (r *random) Uint() uint {
	var num uint
	num = uint(r.Uint64())
	return num
}

func (r *random) Uint8() uint8 {
	var num uint8
	r.read(&num)
	return num
}

func (r *random) Uint16() uint16 {
	var num uint16
	r.read(&num)
	return num
}

func (r *random) Uint32() uint32 {
	var num uint32
	r.read(&num)
	return num
}

func (r *random) Uint64() uint64 {
	var num uint64
	r.read(&num)
	return num
}

func (r *random) Float32() float32 {
	var num float32
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}

func (r *random) Float64() float64 {
	var num float64
	if r.read(&num); num < 0 {
		num = -num
	}
	return num
}
func (r *random) Chars(length int) string {
	return r.CharsWith(length, r.chars)
}

func (r *random) CharsWith(length int, dict []byte) string {
	buf := buff.GetBuff()
	defer buff.PutBuff(buf)
	max := uint8(len(dict))
	for i := 0; i < length; i++ {
		_ = buf.WriteByte(dict[r.Uint8()%max])
	}
	return buf.String()
}

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

func (r *random) read(v any) {
	buf := r.rndBuff(binary.Size(v))
	defer buff.PutBuff(buf)
	_ = binary.Read(buf, binary.BigEndian, v)
}

func (r *random) rndBuff(n int) *bytes.Buffer {
	buf := buff.GetBuff()
	for rnd := r.getBuff(); buf.Len() < n; buf = r.getBuff() {
		_, _ = io.CopyN(buf, rnd, int64(n-buf.Len()))
		r.putBuff(rnd)
	}
	return buf
}

func (r *random) putBuff(buf *bytes.Buffer) {
	if buf.Len() == 0 {
		buff.PutBuff(buf)
	} else {
		r.pool.Put(buf)
	}
}

func (r *random) getBuff() *bytes.Buffer {
	return r.pool.Get().(*bytes.Buffer)
}

func (r *random) newBuff() interface{} {
	buf := buff.GetBuff()
	if err := readBuff(buf, r.size); err != nil {
		panic(err)
	}
	return buf
}
