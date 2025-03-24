package conv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/azeroth-sha/simple"
	"reflect"
	"regexp"
	"strconv"
)

var uintReg = regexp.MustCompile(`^(0[bBoOxX][[:xdigit:]]+|\d+)`)

// ToUint converts any value to uint
func ToUint(v any) uint {
	num, _ := convUint64(v)
	return uint(num)
}

// ToUintE convert any value to uint, return error if failed
func ToUintE(v any) (uint, error) {
	num, err := convUint64(v)
	return uint(num), err
}

// ToUint8 converts any value to uint8
func ToUint8(v any) uint8 {
	num, _ := convUint64(v)
	return uint8(num)
}

// ToUint8E converts any value to uint8, return error if failed
func ToUint8E(v any) (uint8, error) {
	num, err := convUint64(v)
	return uint8(num), err
}

// ToUint16 converts any value to uint16
func ToUint16(v any) uint16 {
	num, _ := convUint64(v)
	return uint16(num)
}

// ToUint16E converts any value to uint16, return error if failed
func ToUint16E(v any) (uint16, error) {
	num, err := convUint64(v)
	return uint16(num), err
}

// ToUint32 converts any value to uint32
func ToUint32(v any) uint32 {
	num, _ := convUint64(v)
	return uint32(num)
}

// ToUint32E converts any value to uint32, return error if failed
func ToUint32E(v any) (uint32, error) {
	num, err := convUint64(v)
	return uint32(num), err
}

// ToUint64 converts any value to uint64
func ToUint64(v any) uint64 {
	num, _ := convUint64(v)
	return num
}

// ToUint64E converts any value to uint64, return error if failed
func ToUint64E(v any) (uint64, error) {
	num, err := convUint64(v)
	return num, err
}

/*
	Package method
*/

func convUint64(v any) (num uint64, err error) {
	switch vv := v.(type) {
	case simple.Int:
		num = uint64(vv.Int())
	case simple.Int8:
		num = uint64(vv.Int8())
	case simple.Int16:
		num = uint64(vv.Int16())
	case simple.Int32:
		num = uint64(vv.Int32())
	case simple.Int64:
		num = uint64(vv.Int64())
	case simple.Uint:
		num = uint64(vv.Uint())
	case simple.Uint8:
		num = uint64(vv.Uint8())
	case simple.Uint16:
		num = uint64(vv.Uint16())
	case simple.Uint32:
		num = uint64(vv.Uint32())
	case simple.Uint64:
		num = vv.Uint64()
	case simple.Float32:
		num = uint64(vv.Float32())
	case simple.Float64:
		num = uint64(vv.Float64())
	case simple.Bool:
		if vv.Bool() {
			num = 1
		}
	case int:
		num = uint64(vv)
	case int8:
		num = uint64(vv)
	case int16:
		num = uint64(vv)
	case int32:
		num = uint64(vv)
	case int64:
		num = uint64(vv)
	case uint:
		num = uint64(vv)
	case uint8:
		num = uint64(vv)
	case uint16:
		num = uint64(vv)
	case uint32:
		num = uint64(vv)
	case uint64:
		num = vv
	case float32:
		num = uint64(vv)
	case float64:
		num = uint64(vv)
	case bool:
		if vv {
			num = 1
		}
	case string:
		num, err = strconv.ParseUint(uintReg.FindString(vv), 0, 0)
	case []byte:
		vvv := string(vv)
		if str := uintReg.FindString(vvv); vvv == str {
			num, err = strconv.ParseUint(str, 0, 0)
		} else {
			bts := vv
			switch len(bts) {
			case 1:
				bts = append(make([]byte, 7), bts...)
				_ = binary.Read(bytes.NewReader(bts), binary.BigEndian, &num)
			case 2:
				bts = append(make([]byte, 6), bts...)
				_ = binary.Read(bytes.NewReader(bts), binary.BigEndian, &num)
			case 4:
				bts = append(make([]byte, 4), bts...)
				_ = binary.Read(bytes.NewReader(bts), binary.BigEndian, &num)
			case 8:
				_ = binary.Read(bytes.NewReader(bts), binary.BigEndian, &num)
			default:
				err = fmt.Errorf("conv: unsupported type %T", v)
			}
		}
	default:
		num, err = convReflectUint64(vv)
	}
	return
}

func convReflectUint64(v any) (uint64, error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convUint64(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convUint64(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convUint64(vv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convUint64(vv.Uint())
	case reflect.Float32, reflect.Float64:
		return convUint64(vv.Float())
	case reflect.String:
		return convUint64(vv.String())
	case reflect.Slice:
		if vv.Type().Elem().Kind() == reflect.Uint8 {
			return convUint64(vv.Bytes())
		}
		fallthrough
	default:
		return 0, fmt.Errorf("conv: unsupported type %T", v)
	}
}
