package odb

import "github.com/azeroth-sha/simple/guid"

// Object 存储对象
type Object interface {
	TableName() (name string)               // 表名
	TableIndex() (fields []string)          // 表索引
	TableField(field string) (value []byte) // 表字段
	TableID() guid.GUID                     // 表主键
	TableNew() (obj Object)                 // 创建新对象
}

// DB 数据库
type DB interface {
	Put(obj Object) (id guid.GUID, err error)
	Get(obj Object, id guid.GUID) error
	Del(obj Object, id guid.GUID) error
	Has(obj Object, index ...string) (has bool, err error)
	Find(obj Object, limit int, index ...string) (arr []guid.GUID, err error)
	Fuzzy(obj Object, limit int, index ...string) (arr []guid.GUID, err error)
	Maintain(obj Object) (err error)
	Close() (err error)
}
