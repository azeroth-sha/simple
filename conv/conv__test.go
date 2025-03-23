package conv

import (
	"fmt"
	"testing"
)

type (
	A int32
	B []byte
)

func TestConvInt(t *testing.T) {
	var (
		aa A = 0x9527
		ab   = fmt.Sprintf(`0x%x`, aa)
		ba B = []byte(fmt.Sprintf(`0x%x`, aa))
		bb   = []byte(fmt.Sprintf(`0x%x`, aa))
	)
	t.Logf("[%d]->(%d)", aa, ToInt(ab))
	t.Logf("[%d]->(%d)", aa, ToInt(aa))
	t.Logf("[%d]->(%d)", aa, ToInt(ba))
	t.Logf("[%d]->(%d)", aa, ToInt(bb))
}
