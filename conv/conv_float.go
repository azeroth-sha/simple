package conv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/azeroth-sha/simple"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var floatReg = regexp.MustCompile(`^(-?\d+(?:\.\d+)?)`)

// ToFloat32 converts any value to float32
func ToFloat32(v any) float32 {
	num, _ := convFloat64(v)
	return float32(num)
}

// ToFloat32E converts any value to float32, return error if failed
func ToFloat32E(v any) (float32, error) {
	num, err := convFloat64(v)
	return float32(num), err
}

// ToFloat64 converts any value to float64
func ToFloat64(v any) float64 {
	num, _ := convFloat64(v)
	return num
}

// ToFloat64E converts any value to float64, return error if failed
func ToFloat64E(v any) (float64, error) {
	return convFloat64(v)
}

/*
  Package method
*/

func convFloat64(v any) (num float64, err error) {
	switch vv := v.(type) {
	case simple.Int:
		num = float64(vv.Int())
	case simple.Int8:
		num = float64(vv.Int8())
	case simple.Int16:
		num = float64(vv.Int16())
	case simple.Int32:
		num = float64(vv.Int32())
	case simple.Int64:
		num = float64(vv.Int64())
	case simple.Uint:
		num = float64(vv.Uint())
	case simple.Uint8:
		num = float64(vv.Uint8())
	case simple.Uint16:
		num = float64(vv.Uint16())
	case simple.Uint32:
		num = float64(vv.Uint32())
	case simple.Uint64:
		num = float64(vv.Uint64())
	case simple.Float32:
		num = float64(vv.Float32())
	case simple.Float64:
		num = vv.Float64()
	case simple.Bool:
		if vv.Bool() {
			num = 1
		}
	case int:
		num = float64(vv)
	case int8:
		num = float64(vv)
	case int16:
		num = float64(vv)
	case int32:
		num = float64(vv)
	case int64:
		num = float64(vv)
	case uint:
		num = float64(vv)
	case uint8:
		num = float64(vv)
	case uint16:
		num = float64(vv)
	case uint32:
		num = float64(vv)
	case uint64:
		num = float64(vv)
	case float32:
		num = float64(vv)
	case float64:
		num = vv
	case bool:
		if vv {
			num = 1
		}
	case string:
		if str := floatReg.FindString(vv); len(str) > 0 && strings.IndexByte(str, '.') > 0 {
			num, err = strconv.ParseFloat(str, 0)
		} else if f, e := ToInt64E(vv); e == nil {
			num = float64(f)
		} else {
			err = e
		}
	case []byte:
		vvv := string(vv)
		if str := floatReg.FindString(vvv); strings.IndexByte(str, '.') > 0 && vvv == str {
			num, err = strconv.ParseFloat(str, 0)
		} else if str = intReg.FindString(vvv); vvv == str {
			if f, e := ToInt64E(str); e == nil {
				num = float64(f)
			} else {
				err = e
			}
		} else {
			switch len(vv) {
			case 4:
				var f32 float32
				_ = binary.Read(bytes.NewReader(vv), binary.BigEndian, &f32)
				num = float64(f32)
			case 8:
				_ = binary.Read(bytes.NewReader(vv), binary.BigEndian, &num)
			default:
				err = fmt.Errorf("conv: unsupported type %T", v)
			}
		}
	default:
		num, err = convReflectFloat64(vv)
	}
	return
}

func convReflectFloat64(v any) (float64, error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convFloat64(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convFloat64(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convFloat64(vv.Int())
	case reflect.Float32, reflect.Float64:
		return convFloat64(vv.Float())
	case reflect.String:
		return convFloat64(vv.String())
	case reflect.Slice:
		if vv.Type().Elem().Kind() == reflect.Uint8 {
			return convFloat64(vv.Bytes())
		}
		fallthrough
	default:
		return 0, fmt.Errorf("conv: unsupported type %T", v)
	}
}
