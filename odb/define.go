package odb

import (
	"errors"
	"github.com/azeroth-sha/simple/guid"
	"github.com/cockroachdb/pebble"
)

var (
	ErrClosed        = pebble.ErrClosed
	ErrNotFound      = pebble.ErrNotFound
	ErrTableNotFound = errors.New(`table not found`)
	ErrIndexNotFound = errors.New(`index not found`)
)

// Object 表结构
type Object interface {
	TableName() (name string)                // 表名
	TableID() (id guid.GUID)                 // 表ID
	TableIndex() (fields []string)           // 表索引
	TableField(field string) (values []byte) // 表字段
}

// Filter 过滤器
type Filter func(index string, value []byte) bool

// Search 查询参数
type Search struct {
	Limit  int      // 查询限制(<=0为不限制)
	Desc   bool     // 是否倒序
	Index  []string // 索引
	Filter Filter   // 过滤器
	UnixL  int64    // 时间范围开始(S)
	UnixU  int64    // 时间范围结束(S)
}

// ODB 数据库接口
type ODB interface {
	// DB 返回pebble.DB
	DB() *pebble.DB
	// Maintain 维护对象的表结构，如果表不存在则创建表
	Maintain(obj Object) (err error)
	// Close 关闭数据库
	Close() (err error)

	// Put 将对象存储到数据库中，并返回对象的唯一标识符
	Put(obj Object) (id guid.GUID, err error)
	// Get 根据唯一标识符从数据库中获取对象
	Get(obj Object, id guid.GUID) (err error)
	// Del 根据唯一标识符从数据库中删除对象
	Del(obj Object, id guid.GUID) (err error)
	// Has 检查对象是否存在于数据库中，可以通过索引加速查找
	Has(obj Object, index ...string) (has bool, err error)
	// Find 根据索引查找对象，并返回对象的唯一标识符列表
	Find(obj Object, search *Search) (all []guid.GUID, err error)
}

/*
  Package define
*/

const (
	keySep = `-`   // 键分隔符
	keyLmt = `.`   // 键限制符
	keyTBL = `tbl` // 表结构
	keyDAT = `dat` // 表数据
	keyIDX = `idx` // 表索引
)

type table struct {
	Name  string   `msgpack:"n"` // 表名
	Index []string `msgpack:"i"` // 表索引
}

func (*table) TableName() string {
	return `table`
}

type inline struct {
	Def *table
	New func() Object
}
