package guid

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"github.com/azeroth-sha/simple/internal"
	"github.com/azeroth-sha/simple/rand"
	"github.com/azeroth-sha/simple/sum"
	"math/big"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	BLen = 12 // GUID字节长度
	SLen = 20 // GUID字符长度
	Base = 36 // GUID转换进制
)

var (
	guidPool = &sync.Pool{New: func() any { return new(big.Int) }}
	endian   = binary.BigEndian
	hID      uint16 // 主机ID
	pID      uint16 // 进程号
	sID      uint32 // 流水号
	NULL     GUID   // 空GUID
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

// HostID returns the host ID.
func (g GUID) HostID() uint16 {
	return binary.BigEndian.Uint16(g[4:6])
}

// ProcessID returns the process ID.
func (g GUID) ProcessID() uint16 {
	return binary.BigEndian.Uint16(g[6:8])
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

// New 生成全局唯一ID
// 4 byte: 时间戳(S)
// 2 byte: 主机号
// 2 byte: 进程号
// 2 byte: 流水号
// 2 byte: 随机数
func New() GUID {
	return NewWithTime(time.Now())
}

// NewWithTime 生成全局唯一ID
func NewWithTime(t time.Time) GUID {
	var id GUID
	endian.PutUint32(id[:4], uint32(t.Unix()))
	endian.PutUint16(id[4:6], hID)
	endian.PutUint16(id[6:8], pID)
	endian.PutUint16(id[8:10], getSerial())
	endian.PutUint16(id[10:], rand.Uint16())
	return id
}

// Parse 解析全局唯一ID
func Parse(id string) (GUID, error) {
	var g GUID
	if err := g.UnmarshalText([]byte(id)); err != nil {
		return g, err
	}
	return g, nil
}

// MustParse 解析全局唯一ID
func MustParse(id string) GUID {
	g, err := Parse(id)
	if err != nil {
		panic(err)
	}
	return g
}

/*
  Package method
*/

func init() {
	hID = getHostID()
	pID = uint16(os.Getpid())
	sID = rand.Uint32()
}

func getHostID() uint16 {
	hid, err := internal.HostID()
	if hid == "" || err != nil {
		hid = rand.Chars(16)
	}
	h := sum.NewCrc16()
	_, _ = h.Write([]byte(hid))
	return h.Sum16()
}

func getSerial() uint16 {
	return uint16(atomic.AddUint32(&sID, 1))
}

func getInt() *big.Int {
	return guidPool.Get().(*big.Int)
}

func putInt(bInt *big.Int) {
	bInt.SetInt64(0)
	guidPool.Put(bInt)
}
