package odb

import (
	"bytes"
	"encoding"
	"github.com/azeroth-sha/simple/guid"
	"github.com/azeroth-sha/simple/internal"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"reflect"
)

func encode(b *bytes.Buffer, v any) error {
	if coding, ok := v.(encoding.BinaryMarshaler); ok {
		if bs, e := coding.MarshalBinary(); e != nil {
			return e
		} else {
			b.Write(bs)
			return nil
		}
	} else {
		return msgpack.NewEncoder(b).Encode(v)
	}
}

func decode(b *bytes.Buffer, v any) error {
	if coding, ok := v.(encoding.BinaryUnmarshaler); ok {
		return coding.UnmarshalBinary(b.Bytes())
	} else {
		return msgpack.NewDecoder(b).Decode(v)
	}
}

func join(buf *bytes.Buffer, sep string, pre []byte, bs ...[]byte) {
	buf.Write(pre)
	for i, b := range bs {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.Write(b)
	}
}

func toBytes(s string) []byte {
	return internal.ToBytes(s)
}

func toString(b []byte) string {
	return internal.ToString(b)
}

func bufReset(bufList ...*bytes.Buffer) {
	for _, buf := range bufList {
		if buf.Len() == 0 {
			continue
		}
		buf.Reset()
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

func mustClose(c io.Closer) {
	if c == nil {
		return
	}
	_ = c.Close()
}

func parseGUID(buf []byte) (id guid.GUID) {
	copy(id[:], buf)
	return
}
