package conv

import (
	"fmt"
	"github.com/azeroth-sha/simple"
	"reflect"
	"strconv"
)

// ToBool converts the value to a bool.
func ToBool(v any) bool {
	b, _ := convBool(v)
	return b
}

// ToBoolE converts the value to a bool, returning an error if it cannot be converted.
func ToBoolE(v any) (b bool, err error) {
	return convBool(v)
}

/*
  Package method
*/

func convBool(v any) (b bool, err error) {
	switch vv := v.(type) {
	case simple.Int:
		b = vv.Int() == 1
	case simple.Int8:
		b = vv.Int8() == 1
	case simple.Int16:
		b = vv.Int16() == 1
	case simple.Int32:
		b = vv.Int32() == 1
	case simple.Int64:
		b = vv.Int64() == 1
	case simple.Uint:
		b = vv.Uint() == 1
	case simple.Uint8:
		b = vv.Uint8() == 1
	case simple.Uint16:
		b = vv.Uint16() == 1
	case simple.Uint32:
		b = vv.Uint32() == 1
	case simple.Uint64:
		b = vv.Uint64() == 1
	case simple.Float32:
		b = vv.Float32() == 1
	case simple.Float64:
		b = vv.Float64() == 1
	case simple.Bool:
		b = vv.Bool()
	case int:
		b = vv == 1
	case int8:
		b = vv == 1
	case int16:
		b = vv == 1
	case int32:
		b = vv == 1
	case int64:
		b = vv == 1
	case uint:
		b = vv == 1
	case uint8:
		b = vv == 1
	case uint16:
		b = vv == 1
	case uint32:
		b = vv == 1
	case uint64:
		b = vv == 1
	case float32:
		b = vv == 1
	case float64:
		b = vv == 1
	case bool:
		b = vv
	case string:
		return strconv.ParseBool(vv)
	case []byte:
		return strconv.ParseBool(string(vv))
	default:
		v, err = convReflectBool(vv)
	}
	return
}

func convReflectBool(v any) (b bool, err error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convBool(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convBool(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convBool(vv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convBool(vv.Uint())
	case reflect.Float32, reflect.Float64:
		return convBool(vv.Float())
	case reflect.String:
		return convBool(vv.String())
	case reflect.Slice:
		if vv.Type().Elem().Kind() == reflect.Uint8 {
			return convBool(vv.Bytes())
		}
		fallthrough
	default:
		return b, fmt.Errorf("conv: unsupported type %T", v)
	}
}
