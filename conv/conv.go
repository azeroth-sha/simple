package conv

import (
	"github.com/azeroth-sha/simple"
	"reflect"
	"strconv"
)

// ToInt conv int
func ToInt(v any) int {
	return int(convInt64(v))
}

// ToInt8 conv int8
func ToInt8(v any) int8 {
	return int8(convInt64(v))
}

// ToInt16 conv int16
func ToInt16(v any) int16 {
	return int16(convInt64(v))
}

// ToInt32 conv int32
func ToInt32(v any) int32 {
	return int32(convInt64(v))
}

// ToInt64 conv int64
func ToInt64(v any) int64 {
	return convInt64(v)
}

// ToUint conv uint
func ToUint(v any) uint {
	return uint(convUint64(v))
}

// ToUint8 conv uint8
func ToUint8(v any) uint8 {
	return uint8(convUint64(v))
}

// ToUint16 conv uint16
func ToUint16(v any) uint16 {
	return uint16(convUint64(v))
}

// ToUint32 conv uint32
func ToUint32(v any) uint32 {
	return uint32(convUint64(v))
}

// ToUint64 conv uint64
func ToUint64(v any) uint64 {
	return convUint64(v)
}

/*
  Package method
*/

func convInt64(v any) int64 {
	switch vv := v.(type) {
	case simple.Int:
		return int64(vv.Int())
	case simple.Int8:
		return int64(vv.Int8())
	case simple.Int16:
		return int64(vv.Int16())
	case simple.Int32:
		return int64(vv.Int32())
	case simple.Int64:
		return vv.Int64()
	case simple.Uint:
		return int64(vv.Uint())
	case simple.Uint8:
		return int64(vv.Uint8())
	case simple.Uint16:
		return int64(vv.Uint16())
	case simple.Uint32:
		return int64(vv.Uint32())
	case simple.Uint64:
		return int64(vv.Uint64())
	case simple.Float32:
		return int64(vv.Float32())
	case simple.Float64:
		return int64(vv.Float64())
	case simple.Bool:
		if vv.Bool() {
			return 1
		} else {
			return 0
		}
	case int:
		return int64(vv)
	case int8:
		return int64(vv)
	case int16:
		return int64(vv)
	case int32:
		return int64(vv)
	case int64:
		return vv
	case uint:
		return int64(vv)
	case uint8:
		return int64(vv)
	case uint16:
		return int64(vv)
	case uint32:
		return int64(vv)
	case uint64:
		return int64(vv)
	case float32:
		return int64(vv)
	case float64:
		return int64(vv)
	case bool:
		if vv {
			return 1
		} else {
			return 0
		}
	case string:
		num, _ := strconv.ParseInt(vv, 0, 0)
		return num
	case []byte:
		num, _ := strconv.ParseInt(string(vv), 0, 0)
		return num
	default:
		return convReflectInt64(vv)
	}
}

func convReflectInt64(v any) int64 {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convInt64(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convInt64(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convInt64(vv.Int())
	case reflect.Float32, reflect.Float64:
		return convInt64(vv.Float())
	case reflect.String:
		return convInt64(vv.String())
	default:
		return 0
	}
}

func convUint64(v any) uint64 {
	switch vv := v.(type) {
	case simple.Int:
		return uint64(vv.Int())
	case simple.Int8:
		return uint64(vv.Int8())
	case simple.Int16:
		return uint64(vv.Int16())
	case simple.Int32:
		return uint64(vv.Int32())
	case simple.Int64:
		return uint64(vv.Int64())
	case simple.Uint:
		return uint64(vv.Uint())
	case simple.Uint8:
		return uint64(vv.Uint8())
	case simple.Uint16:
		return uint64(vv.Uint16())
	case simple.Uint32:
		return uint64(vv.Uint32())
	case simple.Uint64:
		return vv.Uint64()
	case simple.Float32:
		return uint64(vv.Float32())
	case simple.Float64:
		return uint64(vv.Float64())
	case simple.Bool:
		if vv.Bool() {
			return 1
		} else {
			return 0
		}
	case int:
		return uint64(vv)
	case int8:
		return uint64(vv)
	case int16:
		return uint64(vv)
	case int32:
		return uint64(vv)
	case int64:
		return uint64(vv)
	case uint:
		return uint64(vv)
	case uint8:
		return uint64(vv)
	case uint16:
		return uint64(vv)
	case uint32:
		return uint64(vv)
	case uint64:
		return vv
	case float32:
		return uint64(vv)
	case float64:
		return uint64(vv)
	case bool:
		if vv {
			return 1
		} else {
			return 0
		}
	case string:
		num, _ := strconv.ParseUint(vv, 0, 0)
		return num
	case []byte:
		num, _ := strconv.ParseUint(string(vv), 0, 0)
		return num
	default:
		return convReflectUint64(vv)
	}
}

func convReflectUint64(v any) uint64 {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convUint64(vv.Elem().Interface())
	}
	switch vv.Kind() {
	case reflect.Bool:
		return convUint64(vv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convUint64(vv.Int())
	case reflect.Float32, reflect.Float64:
		return convUint64(vv.Float())
	case reflect.String:
		return convUint64(vv.String())
	default:
		return 0
	}
}
