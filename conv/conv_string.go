package conv

import (
	"database/sql"
	"fmt"
	"github.com/azeroth-sha/simple"
	"reflect"
)

const (
	boolTrue  = "true"
	boolFalse = "false"
)

// ToString converts any value to string
func ToString(v any) string {
	str, _ := convString(v)
	return str
}

// ToStringE converts any value to string, return error if failed
func ToStringE(v any) (string, error) {
	return convString(v)
}

/*
  Package method
*/

func convString(v any) (str string, err error) {
	switch vv := v.(type) {
	case sql.NullString:
		str = vv.String
	case simple.Int:
		str = fmt.Sprintf("%d", vv.Int())
	case simple.Int8:
		str = fmt.Sprintf("%d", vv.Int8())
	case simple.Int16:
		str = fmt.Sprintf("%d", vv.Int16())
	case simple.Int32:
		str = fmt.Sprintf("%d", vv.Int32())
	case simple.Int64:
		str = fmt.Sprintf("%d", vv.Int64())
	case simple.Uint:
		str = fmt.Sprintf("%d", vv.Uint())
	case simple.Uint8:
		str = fmt.Sprintf("%d", vv.Uint8())
	case simple.Uint16:
		str = fmt.Sprintf("%d", vv.Uint16())
	case simple.Uint32:
		str = fmt.Sprintf("%d", vv.Uint32())
	case simple.Uint64:
		str = fmt.Sprintf("%d", vv.Uint64())
	case simple.Float32:
		str = fmt.Sprintf("%g", vv.Float32())
	case simple.Float64:
		str = fmt.Sprintf("%g", vv.Float64())
	case simple.Bool:
		if vv.Bool() {
			str = boolTrue
		} else {
			str = boolFalse
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprintf("%d", vv)
	case float32, float64:
		str = fmt.Sprintf("%g", vv)
	case bool:
		if vv {
			str = boolTrue
		} else {
			str = boolFalse
		}
	case string:
		str = vv
	case []byte:
		str = string(vv)
	default:
		str, err = convReflectString(vv)
	}
	return
}

func convReflectString(v any) (string, error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convString(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convString(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convString(vv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convString(vv.Uint())
	case reflect.Float32, reflect.Float64:
		return convString(vv.Float())
	case reflect.String:
		return vv.String(), nil
	case reflect.Slice:
		switch vv.Type().Elem().Kind() {
		case reflect.Uint8:
			return convString(vv.Bytes())
		case reflect.Int32:
			return convString(vv.String())
		default:
			return "", fmt.Errorf("conv: unsupported type %T", v)
		}
	default:
		return "", fmt.Errorf("conv: unsupported type %T", v)
	}
}
