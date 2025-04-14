package conv

import (
	"database/sql"
	"fmt"
	"github.com/azeroth-sha/simple"
	"math"
	"reflect"
	"time"
)

const (
	LayoutDate         = "2006-01-02"
	LayoutTime         = "15:04:05"
	LayoutDateTime     = "2006-01-02 15:04:05"
	LayoutDateTimeMs   = "2006-01-02 15:04:05.999"
	LayoutDateTimeUs   = "2006-01-02 15:04:05.999999"
	LayoutDateTimeNs   = "2006-01-02 15:04:05.999999999"
	LayoutDateTimeZ    = "2006-01-02 15:04:05 Z07:00"
	LayoutDateTimeZMs  = "2006-01-02 15:04:05.999 Z07:00"
	LayoutDateTimeZUs  = "2006-01-02 15:04:05.999999 Z07:00"
	LayoutDateTimeZNs  = "2006-01-02 15:04:05.999999999 Z07:00"
	LayoutDateTimeT    = "2006-01-02T15:04:05"
	LayoutDateTimeZT   = "2006-01-02T15:04:05Z07:00"
	LayoutDateTimeZTMs = "2006-01-02T15:04:05.999Z07:00"
	LayoutDateTimeZTUs = "2006-01-02T15:04:05.999999Z07:00"
	LayoutDateTimeZTNs = "2006-01-02T15:04:05.999999999Z07:00"
)

// ToTime converts any value to time.Time.
func ToTime(v any, loc ...*time.Location) time.Time {
	tm, _ := ToTimeE(v, loc...)
	return tm
}

// ToTimeE converts any value to time.Time, returning an error if the conversion fails.
func ToTimeE(v any, loc ...*time.Location) (time.Time, error) {
	return convTime(v, loc...)
}

// ParseTime parses a string into a time.Time value.
func ParseTime(s string, loc ...*time.Location) time.Time {
	tm, _ := ParseTimeE(s, loc...)
	return tm
}

// ParseTimeE parses a string into a time.Time value, returning an error if the conversion fails.
func ParseTimeE(s string, loc ...*time.Location) (tm time.Time, err error) {
	l := time.Local
	if len(loc) > 0 {
		l = loc[0]
	}
	for _, f := range layouts {
		if f.z {
			if tm, err = time.Parse(f.f, s); err == nil {
				tm.In(l)
				return tm, nil
			}
		} else if tm, err = time.ParseInLocation(f.f, s, l); err == nil {
			return tm, nil
		}
	}
	return tm, fmt.Errorf("conv: unable to parse date: %s", s)
}

// FormatTime formats a time.Time value into a string.
func FormatTime(v any, l ...string) string {
	str, _ := FormatTimeE(v, l...)
	return str
}

// FormatTimeE formats a time.Time value into a string, returning an error if the conversion fails.
func FormatTimeE(v any, l ...string) (string, error) {
	if tm, err := convTime(v, nil); err != nil {
		return "", err
	} else if !tm.IsZero() {
		f := LayoutDateTimeZT
		if len(l) > 0 {
			f = l[0]
		}
		return tm.Format(f), nil
	} else {
		return "", fmt.Errorf("conv: time is zero: %v", v)
	}
}

// Unix returns the Unix time of the given value.
func Unix(v any) int64 {
	unixSec, _ := UnixE(v)
	return unixSec
}

// UnixE returns the Unix time of the given value, returning an error if the conversion fails.
func UnixE(v any) (int64, error) {
	tm, err := ToTimeE(v)
	if err != nil {
		return 0, err
	} else if tm.IsZero() {
		return 0, nil
	}
	return tm.Unix(), nil
}

// UnixMs returns the Unix millisecond time of the given value.
func UnixMs(v any) int64 {
	unixMS, _ := UnixMsE(v)
	return unixMS
}

// UnixMsE returns the Unix millisecond time of the given value, returning an error if the conversion fails.
func UnixMsE(v any) (int64, error) {
	tm, err := ToTimeE(v)
	if err != nil {
		return 0, err
	} else if tm.IsZero() {
		return 0, nil
	}
	return tm.UnixMilli(), nil
}

/*
	Package method
*/

type format struct {
	f string // format string
	z bool   // has zone info
}

var layouts = []*format{
	{LayoutDate, false},
	{LayoutTime, false},
	{LayoutDateTime, false},
	{LayoutDateTimeMs, false},
	{LayoutDateTimeUs, false},
	{LayoutDateTimeNs, false},
	{LayoutDateTimeZ, true},
	{LayoutDateTimeZMs, true},
	{LayoutDateTimeZUs, true},
	{LayoutDateTimeZNs, true},
	{LayoutDateTimeT, false},
	{LayoutDateTimeZT, true},
	{LayoutDateTimeZTMs, true},
	{LayoutDateTimeZTUs, true},
	{LayoutDateTimeZTNs, true},
}

func convTime(v any, loc ...*time.Location) (tm time.Time, err error) {
	switch vv := v.(type) {
	case sql.NullTime:
		tm = vv.Time
	case time.Time:
		tm = vv
	case simple.Int:
		tm, err = convTime(int64(vv.Int()), loc...)
	case simple.Int32:
		tm, err = convTime(int64(vv.Int32()), loc...)
	case simple.Int64:
		tm, err = convTime(vv.Int64(), loc...)
	case simple.Uint:
		tm, err = convTime(int64(vv.Uint()), loc...)
	case simple.Uint32:
		tm, err = convTime(int64(vv.Uint32()), loc...)
	case simple.Uint64:
		tm, err = convTime(int64(vv.Uint64()), loc...)
	case int:
		tm, err = convTime(int64(vv), loc...)
	case int32:
		tm, err = convTime(int64(vv), loc...)
	case int64:
		l := time.Local
		if len(loc) > 0 {
			l = loc[0]
		}
		var sec, nsec int64 = vv, 0
		for vv > math.MaxUint32 {
			sec = vv / 1000
			nsec = (vv % 1000) * 1000000
			vv /= 1000
		}
		tm = time.Unix(sec, nsec).In(l)
	case uint:
		tm, err = convTime(int64(vv), loc...)
	case uint32:
		tm, err = convTime(int64(vv), loc...)
	case uint64:
		tm, err = convTime(int64(vv), loc...)
	case string:
		tm, err = ParseTimeE(vv, loc...)
	default:
		tm, err = convReflectTime(v, loc...)
	}
	return
}

func convReflectTime(v any, loc ...*time.Location) (time.Time, error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return convTime(vv.Elem().Interface(), loc...)
	}
	switch vv.Kind() {
	case reflect.Int32, reflect.Int64:
		return convTime(vv.Int(), loc...)
	case reflect.Uint32, reflect.Uint64:
		return convTime(vv.Uint(), loc...)
	case reflect.String:
		return convTime(vv.String(), loc...)
	default:
		return time.Time{}, fmt.Errorf("conv: unsupported type %T", v)
	}
}
