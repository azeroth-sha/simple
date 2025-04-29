package guid

import (
	"github.com/azeroth-sha/simple/internal"
	"github.com/azeroth-sha/simple/rand"
	"hash/fnv"
	"os"
	"time"
)

var (
	NULL    GUID // 空GUID
	hostGen Generator
)

// New returns a new GUID.
func New() GUID {
	return NewWithTime(time.Now())
}

// NewWithTime returns a new GUID with the given time.
func NewWithTime(tm time.Time) GUID {
	return hostGen.NewWithTime(tm)
}

// Parse parses a GUID from a string.
func Parse(s string) (GUID, error) {
	var id GUID
	return id, id.UnmarshalText([]byte(s))
}

// MustParse parses a GUID from a string.
func MustParse(s string) GUID {
	id, _ := Parse(s)
	return id
}

// NewGenerator returns a new Generator.
func NewGenerator(mark uint32) Generator {
	return &adapter{
		mark:   mark,
		serial: rand.Uint32(),
	}
}

/*
  Package Private functions
*/

func init() {
	hID := getHostID()
	pID := uint32(os.Getpid())
	hostGen = NewGenerator(hID<<16 | pID)
}

func getHostID() uint32 {
	hid, err := internal.HostID()
	if hid == "" || err != nil {
		hid = rand.Chars(16)
	}
	h := fnv.New32()
	_, _ = h.Write([]byte(hid))
	return h.Sum32()
}
