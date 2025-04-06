package odb

import (
	"bytes"
	"github.com/azeroth-sha/simple/guid"
	"github.com/azeroth-sha/simple/internal"
	"github.com/ugorji/go/codec"
	"io"
)

func joinKey(pre []byte, bs ...[]byte) []byte {
	if len(pre) == 0 {
		return bytes.Join(bs, toBs(joinChar))
	} else {
		return append(pre, bytes.Join(bs, toBs(joinChar))...)
	}
}

func objDATKey(keyBuf *bytes.Buffer, pre []byte, obj Object, id guid.GUID) {
	_, _ = keyBuf.Write(pre)
	_, _ = keyBuf.Write(toBs(prefixDat))
	_, _ = keyBuf.Write(toBs(joinChar))
	_, _ = keyBuf.Write(toBs(obj.TableName()))
	_, _ = keyBuf.Write(toBs(joinChar))
	_, _ = keyBuf.Write(id.Bytes())
}

func objIDXKey(keyBuf *bytes.Buffer, pre []byte, obj Object, index string, id guid.GUID) {
	_, _ = keyBuf.Write(pre)
	_, _ = keyBuf.Write(toBs(prefixIdx))
	_, _ = keyBuf.Write(toBs(joinChar))
	_, _ = keyBuf.Write(toBs(obj.TableName()))
	_, _ = keyBuf.Write(toBs(joinChar))
	_, _ = keyBuf.Write(toBs(index))
	_, _ = keyBuf.Write(toBs(joinChar))
	_, _ = keyBuf.Write(obj.TableField(index))
	_, _ = keyBuf.Write(toBs(joinChar))
	_, _ = keyBuf.Write(id.Bytes())
}

func encode(w io.Writer, v interface{}) error {
	mh := new(codec.MsgpackHandle)
	return codec.NewEncoder(w, mh).Encode(v)
}

func decode(r io.Reader, v interface{}) error {
	mh := new(codec.MsgpackHandle)
	return codec.NewDecoder(r, mh).Decode(v)
}

func toBs(s string) []byte {
	return internal.ToBytes(s)
}
