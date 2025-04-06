package odb

import "strings"

const (
	joinChar  = `-`   // 分割符
	limitChar = `.`   // 限制符
	prefixDef = `def` // 定义前缀
	prefixDat = `dat` // 数据前缀
	prefixIdx = `idx` // 索引前缀
)

type define struct {
	Name  string   `msg:"name"`
	Index []string `msg:"index"`
}

func (d *define) TableKey() (_ string) {
	b := new(strings.Builder)
	b.WriteString(prefixDef)
	b.WriteString(joinChar)
	b.WriteString(d.Name)
	return b.String()
}
