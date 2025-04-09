package rand

import (
	"crypto/rand"
	"io"
	"sync"
)

var (
	_reader     = rand.Reader                 // 加密安全的随机数生成器
	_readMu     = new(sync.Mutex)             // 用于保护随机数读取的互斥锁
	defaultRand = NewRandom(256, DefaultDict) // 默认的随机数生成器实例
)

// Int 生成一个随机的 int 类型整数
func Int() int {
	return defaultRand.Int()
}

// Int8 生成一个随机的 int8 类型整数
func Int8() int8 {
	return defaultRand.Int8()
}

// Int16 生成一个随机的 int16 类型整数
func Int16() int16 {
	return defaultRand.Int16()
}

// Int32 生成一个随机的 int32 类型整数
func Int32() int32 {
	return defaultRand.Int32()
}

// Int64 生成一个随机的 int64 类型整数
func Int64() int64 {
	return defaultRand.Int64()
}

// Uint 生成一个随机的 uint 类型整数
func Uint() uint {
	return defaultRand.Uint()
}

// Uint8 生成一个随机的 uint8 类型整数
func Uint8() uint8 {
	return defaultRand.Uint8()
}

// Uint16 生成一个随机的 uint16 类型整数
func Uint16() uint16 {
	return defaultRand.Uint16()
}

// Uint32 生成一个随机的 uint32 类型整数
func Uint32() uint32 {
	return defaultRand.Uint32()
}

// Uint64 生成一个随机的 uint64 类型整数
func Uint64() uint64 {
	return defaultRand.Uint64()
}

// Float32 生成一个随机的 float32 类型浮点数
func Float32() float32 {
	return defaultRand.Float32()
}

// Float64 生成一个随机的 float64 类型浮点数
func Float64() float64 {
	return defaultRand.Float64()
}

// Chars 生成指定长度的随机字符串，使用默认字符集
func Chars(length int) string {
	return defaultRand.Chars(length)
}

// CharsWith 生成指定长度的随机字符串，使用自定义字符集
func CharsWith(length int, dict []byte) string {
	return defaultRand.CharsWith(length, dict)
}

// TextWith 生成指定长度的随机文本，使用自定义字符集
func TextWith(length int, dict []rune) string {
	return defaultRand.TextWith(length, dict)
}

// NewRandom 创建一个新的随机数生成器实例
func NewRandom(size int, chars string) Random {
	r := &random{
		pool:  new(sync.Pool), // 用于管理缓冲区的 sync.Pool
		size:  size,           // 缓冲区大小
		chars: []byte(chars),  // 默认字符集
	}
	r.pool.New = r.newBuff // 设置缓冲区的创建函数
	return r
}

/*
  Package method
*/

// readBuff 从加密安全的随机数生成器中读取指定长度的随机字节到缓冲区
func readBuff(w io.Writer, n int) error {
	_readMu.Lock()
	defer _readMu.Unlock()
	_, err := io.CopyN(w, _reader, int64(n))
	return err
}
