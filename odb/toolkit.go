package odb

import (
	"bytes"
	"encoding"
	"github.com/azeroth-sha/simple/buff"
	"github.com/azeroth-sha/simple/guid"
	"github.com/azeroth-sha/simple/internal"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

func encode(b *bytes.Buffer, v any) error {
	if coding, ok := v.(encoding.BinaryMarshaler); ok {
		if bs, e := coding.MarshalBinary(); e != nil {
			return e
		} else {
			b.Write(bs)
		}
	} else {
		return msgpack.NewEncoder(b).Encode(v)
	}
	return nil
}

func decode(b *bytes.Buffer, v any) error {
	if coding, ok := v.(encoding.BinaryUnmarshaler); ok {
		return coding.UnmarshalBinary(b.Bytes())
	} else {
		return msgpack.NewDecoder(b).Decode(v)
	}
}

func join(buf *bytes.Buffer, pre []byte, sep string, bs ...[]byte) {
	if len(pre) > 0 {
		buf.Write(pre)
	}
	for i, b := range bs {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.Write(b)
	}
}

func toBts(s string) []byte {
	return internal.ToBytes(s)
}

func toStr(b []byte) string {
	return internal.ToString(b)
}

func getBuf() *bytes.Buffer {
	return buff.GetBuff()
}

func putBuf(all ...*bytes.Buffer) {
	for _, b := range all {
		buff.PutBuff(b)
	}
}

func resetBuf(bufList ...*bytes.Buffer) {
	for _, buf := range bufList {
		if buf.Len() > 0 {
			buf.Reset()
		}
	}
}

func reflectNew(o Object) func() Object {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return func() Object {
		return reflect.New(t).Interface().(Object)
	}
}

func discardErr(f func() error) {
	if f == nil {
		return
	}
	_ = f()
}

func toGUID(buf []byte) (id guid.GUID) {
	if len(buf) == guid.BLen {
		copy(id[:], buf)
	}
	return
}
