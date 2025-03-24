package tools

import (
	"bytes"
	"github.com/azeroth-sha/simple/internal"
	"strings"
)

// TrimSpace return the string without space
func TrimSpace(str string) string {
	bs := internal.ToBytes(str)
	if i := bytes.IndexByte(bs, 0); i >= 0 {
		str = str[:i]
	}
	return strings.TrimSpace(str)
}

// Split return the string split by sep
func Split(str, sep string) []string {
	all := strings.Split(str, sep)
	for i, v := range all {
		all[i] = TrimSpace(v)
	}
	return all
}
