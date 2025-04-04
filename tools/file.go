package tools

import (
	"github.com/azeroth-sha/simple/buff"
	"io"
	"net/http"
	"os"
	"time"
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

// FileTemp creates a temporary file
func FileTemp(r io.Reader, n string, d time.Duration) (string, error) {
	fd, err := os.CreateTemp(os.TempDir(), n)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fd.Close()
	}()
	if _, err = io.Copy(fd, r); err != nil {
		return "", err
	}
	if d > 0 {
		_ = time.AfterFunc(d, func() {
			_ = os.Remove(fd.Name())
		})
	}
	return fd.Name(), nil
}
