package guid

import (
	"encoding/binary"
	"github.com/azeroth-sha/simple/rand"
	"sync/atomic"
	"time"
)

var endian = binary.BigEndian

type Generator interface {
	New() GUID
	NewWithTime(tm time.Time) GUID
}

type adapter struct {
	mark   uint32
	serial uint32
}

func (a *adapter) New() GUID {
	return a.NewWithTime(time.Now())
}

func (a *adapter) NewWithTime(tm time.Time) (id GUID) {
	endian.PutUint32(id[:4], uint32(tm.Unix()))
	endian.PutUint32(id[4:8], a.mark)
	endian.PutUint16(id[8:10], a.getSerial())
	endian.PutUint16(id[10:], rand.Uint16())
	return id
}

/*
  Package Private functions
*/

func (a *adapter) getSerial() uint16 {
	n := atomic.AddUint32(&a.serial, 1)
	return uint16(n)
}
