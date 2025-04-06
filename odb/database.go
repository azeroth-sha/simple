package odb

import (
	"bytes"
	"errors"
	"github.com/azeroth-sha/simple/buff"
	"github.com/azeroth-sha/simple/guid"
	"github.com/cockroachdb/pebble"
	"golang.org/x/exp/slices"
	"sync"
	"sync/atomic"
)

type database struct {
	closed *atomic.Bool
	mutex  *sync.RWMutex
	prefix []byte
	db     *pebble.DB
	tables map[string]*define
}

func (d *database) Put(obj Object) (id guid.GUID, _ error) {
	if d.isClosed() {
		return id, pebble.ErrClosed
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	// TODO
	return id, nil
}

func (d *database) Get(obj Object, id guid.GUID) error {
	//TODO implement me
	panic("implement me")
}

func (d *database) Del(obj Object, id guid.GUID) error {
	//TODO implement me
	panic("implement me")
}

func (d *database) Has(obj Object, index ...string) (has bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *database) Find(obj Object, limit int, index ...string) (arr []guid.GUID, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *database) Fuzzy(obj Object, limit int, index ...string) (arr []guid.GUID, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *database) Maintain(obj Object) (_ error) {
	if d.isClosed() {
		return pebble.ErrClosed
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.maintain(obj)
}

func (d *database) Close() (err error) {
	if d.isClosed() {
		return nil
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.closed.Store(true)
	return d.db.Close()
}

/*
  Package method
*/

func (d *database) put(obj Object) (id guid.GUID, _ error) {
	key := buff.GetBuff()
	defer buff.PutBuff(key)
	val := buff.GetBuff()
	defer buff.PutBuff(val)
	batch := d.db.NewBatch()
	defer batch.Close()
	id = obj.TableID()

	objDATKey(key, d.prefix, obj, id)
	if e := encode(val, obj); e != nil {
		return id, e
	} else if e = batch.Set(key.Bytes(), val.Bytes(), nil); e != nil {
		return id, e
	}
}
func (d *database) get(obj Object, id guid.GUID) error {
	key := buff.GetBuff()
	defer buff.PutBuff(key)
	objDATKey(key, d.prefix, obj, id)
	if v, c, e := d.db.Get(key.Bytes()); c != nil {
		return e
	} else {
		defer c.Close()
		return decode(bytes.NewBuffer(v), obj)
	}
}
func (d *database) del(obj Object, id guid.GUID) error {
	key := joinKey(d.prefix, toBs(prefixDat), toBs(obj.TableName()), id.Bytes())
	return d.db.Delete(key, pebble.Sync)
}

func (d *database) has(obj Object, index string) (has bool, err error) {
	key := joinKey(d.prefix, toBs(prefixIdx), toBs(obj.TableName()), obj.TableID().Bytes(), toBs(index))
}
func (d *database) find(obj Object, limit int, index ...string) (arr []guid.GUID, err error)  {}
func (d *database) fuzzy(obj Object, limit int, index ...string) (arr []guid.GUID, err error) {}

func (d *database) maintain(obj Object) error {
	newDef := &define{
		Index: obj.TableIndex(),
		Name:  obj.TableName(),
	}
	key := joinKey(d.prefix, []byte(newDef.TableKey()))
	newVal := buff.GetBuff()
	defer buff.PutBuff(newVal)
	if e := encode(newVal, newDef); e != nil {
		return e
	}
	if v, c, e := d.db.Get(key); e != nil && !errors.Is(e, pebble.ErrNotFound) {
		return e
	} else if e == nil {
		oldVal := buff.GetBuff()
		defer buff.PutBuff(oldVal)
		oldVal.Write(v)
		oldDef := new(define)
		if e = c.Close(); e != nil {
			return e
		} else if e = decode(oldVal, oldDef); e != nil {
			return e
		} else if e = d.resetIdx(obj, newDef, oldDef); e != nil {
			return e
		}
	}
	if e := d.db.Set(key, newVal.Bytes(), pebble.Sync); e != nil {
		return e
	} else {
		d.tables[newDef.Name] = newDef
		return nil
	}
}

func (d *database) resetIdx(obj Object, newDef, oldDef *define) error {
	for i := 0; i < len(newDef.Index); i++ {
		index := newDef.Index[i]
		if slices.Contains(oldDef.Index, index) {
			continue
		} else if e := d.addIdx(obj, index); e != nil {
			return e
		}
	}
	for i := 0; i < len(oldDef.Index); i++ {
		index := oldDef.Index[i]
		if slices.Contains(newDef.Index, index) {
			continue
		} else if e := d.delIdx(obj, index); e != nil {
			return e
		}
	}
	return nil
}

func (d *database) delIdx(obj Object, index string) error {
	sk := joinKey(d.prefix, toBs(prefixIdx), toBs(obj.TableName()), toBs(index+joinChar))
	ek := joinKey(d.prefix, toBs(prefixIdx), toBs(obj.TableName()), toBs(index+limitChar))
	return d.db.DeleteRange(sk, ek, pebble.Sync)
}

func (d *database) addIdx(obj Object, index string) error {
	snap := d.db.NewSnapshot()
	defer snap.Close()
	iter, err := snap.NewIter(&pebble.IterOptions{
		LowerBound: joinKey(d.prefix, toBs(prefixDat), toBs(obj.TableName()+joinChar)),
		UpperBound: joinKey(d.prefix, toBs(prefixDat), toBs(obj.TableName()+limitChar)),
	})
	if err != nil {
		return err
	}
	defer iter.Close()
	batch := d.db.NewBatch()
	defer batch.Close()
EXIT:
	for iter.First(); iter.Valid(); iter.Next() {
		o := obj.TableNew()
		if v, e := iter.ValueAndErr(); e != nil {
			err = e
			break EXIT
		} else if err = decode(bytes.NewBuffer(v), o); err != nil {
			break EXIT
		}
		key := joinKey(d.prefix, toBs(prefixIdx), toBs(obj.TableName()), toBs(index), o.TableID().Bytes(), obj.TableField(index))
		if err = batch.Set(key, o.TableID().Bytes(), nil); err != nil {
			break EXIT
		} else if batch.Len() > 256 {
			if err = batch.Commit(pebble.Sync); err != nil {
				break EXIT
			} else {
				batch.Reset()
			}
		}

	}

	if err == nil && batch.Len() > 0 {
		err = batch.Commit(pebble.Sync)
	}
	return err
}

func (d *database) isClosed() bool {
	return d.closed.Load()
}
