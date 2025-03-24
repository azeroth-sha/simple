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

var intReg = regexp.MustCompile(`^(-?0[bBoOxX][[:xdigit:]]+|-?\d+)`)

// ToInt converts any value to int
func ToInt(v any) int {
	num, _ := convInt64(v)
	return int(num)
}

// ToIntE convert any value to int, return error if failed
func ToIntE(v any) (int, error) {
	num, err := convInt64(v)
	return int(num), err
}

// ToInt8 converts any value to int8
func ToInt8(v any) int8 {
	num, _ := convInt64(v)
	return int8(num)
}

// ToInt8E converts any value to int8, return error if failed
func ToInt8E(v any) (int8, error) {
	num, err := convInt64(v)
	return int8(num), err
}

// ToInt16 converts any value to int16
func ToInt16(v any) int16 {
	num, _ := convInt64(v)
	return int16(num)
}

// ToInt16E converts any value to int16, return error if failed
func ToInt16E(v any) (int16, error) {
	num, err := convInt64(v)
	return int16(num), err
}

// ToInt32 converts any value to int32
func ToInt32(v any) int32 {
	num, _ := convInt64(v)
	return int32(num)
}

// ToInt32E converts any value to int32, return error if failed
func ToInt32E(v any) (int32, error) {
	num, err := convInt64(v)
	return int32(num), err
}

// ToInt64 converts any value to int64
func ToInt64(v any) int64 {
	num, _ := convInt64(v)
	return num
}

// ToInt64E converts any value to int64, return error if failed
func ToInt64E(v any) (int64, error) {
	num, err := convInt64(v)
	return num, err
}

/*
  Package method
*/

func convInt64(v any) (num int64, err error) {
	switch vv := v.(type) {
	case simple.Int:
		num = int64(vv.Int())
	case simple.Int8:
		num = int64(vv.Int8())
	case simple.Int16:
		num = int64(vv.Int16())
	case simple.Int32:
		num = int64(vv.Int32())
	case simple.Int64:
		num = vv.Int64()
	case simple.Uint:
		num = int64(vv.Uint())
	case simple.Uint8:
		num = int64(vv.Uint8())
	case simple.Uint16:
		num = int64(vv.Uint16())
	case simple.Uint32:
		num = int64(vv.Uint32())
	case simple.Uint64:
		num = int64(vv.Uint64())
	case simple.Float32:
		num = int64(vv.Float32())
	case simple.Float64:
		num = int64(vv.Float64())
	case simple.Bool:
		if vv.Bool() {
			num = 1
		}
	case int:
		num = int64(vv)
	case int8:
		num = int64(vv)
	case int16:
		num = int64(vv)
	case int32:
		num = int64(vv)
	case int64:
		num = vv
	case uint:
		num = int64(vv)
	case uint8:
		num = int64(vv)
	case uint16:
		num = int64(vv)
	case uint32:
		num = int64(vv)
	case uint64:
		num = int64(vv)
	case float32:
		num = int64(vv)
	case float64:
		num = int64(vv)
	case bool:
		if vv {
			num = 1
		}
	case string:
		num, err = strconv.ParseInt(intReg.FindString(vv), 0, 0)
	case []byte:
		vvv := string(vv)
		if str := intReg.FindString(vvv); vvv == str {
			num, _ = strconv.ParseInt(str, 0, 0)
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
		num, err = convReflectInt64(vv)
	}
	return
}

func convReflectInt64(v any) (int64, error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convInt64(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convInt64(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convInt64(vv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convInt64(vv.Uint())
	case reflect.Float32, reflect.Float64:
		return convInt64(vv.Float())
	case reflect.String:
		return convInt64(vv.String())
	case reflect.Slice:
		if vv.Type().Elem().Kind() == reflect.Uint8 {
			return convInt64(vv.Bytes())
		}
		fallthrough
	default:
		return 0, fmt.Errorf("conv: unsupported type %T", v)
	}
}
