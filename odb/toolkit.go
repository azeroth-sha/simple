package odb

import (
	"bytes"
	"github.com/azeroth-sha/simple/guid"
	"github.com/azeroth-sha/simple/internal"
	"github.com/ugorji/go/codec"
	"io"
	"reflect"
)

func encode(b *bytes.Buffer, v any) error {
	h := new(codec.MsgpackHandle)
	return codec.NewEncoder(b, h).Encode(v)
}

func decode(b *bytes.Buffer, v any) error {
	h := new(codec.MsgpackHandle)
	return codec.NewDecoder(b, h).Decode(v)
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
	if t.Kind() != reflect.Ptr {
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

func idIntersect(a, b []guid.GUID) []guid.GUID {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	m := make(map[guid.GUID]struct{})
	for _, id := range a {
		if _, ok := m[id]; ok {
			continue
		}
		m[id] = struct{}{}
	}
	a = a[:0]
	l := len(b)
	for i := 0; i < l; i++ {
		if _, ok := m[b[i]]; ok {
			a = append(a, b[i])
		}
	}
	return a
}
