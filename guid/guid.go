package guid

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

const (
	BLen = 12 // GUID字节长度
	SLen = 20 // GUID字符长度
	Base = 36 // GUID转换进制
)

// GUID is the global unique ID.
type GUID [BLen]byte

func (g GUID) Reset() {
	for i := 0; i < BLen; i++ {
		g[i] = 0
	}
}

// Empty returns true if the ID is empty.
func (g GUID) Empty() bool {
	return g == NULL
}

// String returns the string.
func (g GUID) String() string {
	bInt := getInt()
	defer putInt(bInt)
	bInt.SetBytes(g[:])
	id := bInt.Text(Base)
	if len(id) < SLen {
		id = strings.Repeat("0", SLen-len(id)) + id
	}
	return id
}

// Bytes returns the byte slice.
func (g GUID) Bytes() []byte {
	return g[:]
}

// Unix returns the timestamp.
func (g GUID) Unix() int64 {
	return int64(binary.BigEndian.Uint32(g[:4]))
}

// MarkID returns the mark ID.
func (g GUID) MarkID() uint32 {
	return binary.BigEndian.Uint32(g[4:8])
}

// Serial returns the serial number.
func (g GUID) Serial() uint16 {
	return binary.BigEndian.Uint16(g[8:10])
}

// Random returns the random number.
func (g GUID) Random() uint16 {
	return binary.BigEndian.Uint16(g[10:])
}

// Equal returns true if the two IDs are equal.
func (g GUID) Equal(id GUID) bool {
	return g == id
}

// Lt returns true if the first ID is less than the second ID.
func (g GUID) Lt(id GUID) bool {
	return bytes.Compare(g[:], id[:]) < 0
}

// Lte returns true if the first ID is less than or equal to the second ID.
func (g GUID) Lte(id GUID) bool {
	return bytes.Compare(g[:], id[:]) <= 0
}

// Gt returns true if the first ID is greater than the second ID.
func (g GUID) Gt(id GUID) bool {
	return bytes.Compare(g[:], id[:]) > 0
}

// Gte returns true if the first ID is greater than or equal to the second ID.
func (g GUID) Gte(id GUID) bool {
	return bytes.Compare(g[:], id[:]) >= 0
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (g GUID) MarshalBinary() (data []byte, err error) {
	return g[:], nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (g *GUID) UnmarshalBinary(data []byte) error {
	switch len(data) {
	case 0:
		g.Reset()
	case BLen:
		_ = copy(g[:], data)
	case SLen:
		return g.UnmarshalText(data)
	default:
		return fmt.Errorf("%s: %v", os.ErrInvalid, data)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (g GUID) MarshalText() (text []byte, err error) {
	return []byte(g.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (g *GUID) UnmarshalText(text []byte) error {
	switch len(text) {
	case 0, 4:
		g.Reset()
	case BLen:
		return g.UnmarshalBinary(text)
	case SLen:
		bInt := getInt()
		defer putInt(bInt)
		if _, ok := bInt.SetString(string(text), Base); !ok {
			return fmt.Errorf("%s: %v", os.ErrInvalid, text)
		} else {
			copy(g[:], bInt.Bytes())
		}
	default:
		return fmt.Errorf("%s: %s", os.ErrInvalid, text)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (g GUID) MarshalJSON() ([]byte, error) {
	out := make([]byte, SLen+2)
	out[0] = '"'
	out[SLen+1] = '"'
	buf, _ := g.MarshalText()
	copy(out[1:], buf)
	return out, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (g *GUID) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, `"`)
	return g.UnmarshalText(data)
}

// Value implements the driver.Valuer interface.
func (g GUID) Value() (driver.Value, error) {
	return g.String(), nil
}

// Scan implements the sql.Scanner interface.
func (g *GUID) Scan(src any) error {
	switch v := src.(type) {
	case string:
		return g.UnmarshalText([]byte(v))
	case []byte:
		return g.UnmarshalBinary(v)
	case nil:
		g.Reset()
	default:
		return fmt.Errorf("%s: %v", os.ErrInvalid, src)
	}
	return nil
}
