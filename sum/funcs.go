package sum

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm3"
)

// CRC16 计算并返回输入字节切片的CRC16校验值。
func CRC16(bs []byte) []byte {
	h := NewCrc16()
	_, _ = h.Write(bs)
	return h.Sum(nil)
}

// CRC16Hex 计算输入字节切片的CRC16校验值，并返回其十六进制字符串表示。
func CRC16Hex(bs []byte) string {
	return hex.EncodeToString(CRC16(bs))
}

// MD5 计算并返回输入字节切片的MD5哈希值。
func MD5(bs []byte) []byte {
	h := md5.New()
	_, _ = h.Write(bs)
	return h.Sum(nil)
}

// MD5Hex 计算输入字节切片的MD5哈希值，并返回其十六进制字符串表示。
func MD5Hex(bs []byte) string {
	return hex.EncodeToString(MD5(bs))
}

// SM3 计算并返回输入字节切片的SM3哈希值。
func SM3(bs []byte) []byte {
	h := sm3.New()
	_, _ = h.Write(bs)
	return h.Sum(nil)
}

// SM3Hex 计算输入字节切片的SM3哈希值，并返回其十六进制字符串表示。
func SM3Hex(bs []byte) string {
	return hex.EncodeToString(SM3(bs))
}
