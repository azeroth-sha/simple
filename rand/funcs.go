package rand

import (
	"crypto/rand"
	"io"
	"sync"
)

var (
	_reader     = rand.Reader
	_readMu     = new(sync.Mutex)
	defaultRand = NewRandom(256, DefaultDict)
)

func Int() int {
	return defaultRand.Int()
}

func Int8() int8 {
	return defaultRand.Int8()
}

func Int16() int16 {
	return defaultRand.Int16()
}

func Int32() int32 {
	return defaultRand.Int32()
}

func Int64() int64 {
	return defaultRand.Int64()
}

func Uint() uint {
	return defaultRand.Uint()
}

func Uint8() uint8 {
	return defaultRand.Uint8()
}

func Uint16() uint16 {
	return defaultRand.Uint16()
}

func Uint32() uint32 {
	return defaultRand.Uint32()
}

func Uint64() uint64 {
	return defaultRand.Uint64()
}

func Float32() float32 {
	return defaultRand.Float32()
}

func Float64() float64 {
	return defaultRand.Float64()
}

func Chars(length int) string {
	return defaultRand.Chars(length)
}

func CharsWith(length int, dict []byte) string {
	return defaultRand.CharsWith(length, dict)
}

func TextWith(length int, dict []rune) string {
	return defaultRand.TextWith(length, dict)
}

func NewRandom(size int, chars string) Random {
	r := &random{
		pool:  new(sync.Pool),
		size:  size,
		chars: []byte(chars),
	}
	r.pool.New = r.newBuff
	return r
}

/*
  Package method
*/

func readBuff(w io.Writer, n int) error {
	_readMu.Lock()
	defer _readMu.Unlock()
	_, err := io.CopyN(w, _reader, int64(n))
	return err
}
