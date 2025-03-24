package tools

import (
	"github.com/azeroth-sha/simple/buff"
	"io"
	"net/http"
	"os"
)

// FileType returns the file type
func FileType(name string) string {
	buf := buff.GetBuff()
	defer buff.PutBuff(buf)
	if fd, err := os.Open(name); err == nil {
		_, _ = io.CopyN(buf, fd, 512)
		_ = fd.Close()
	}
	return http.DetectContentType(buf.Bytes())
}
