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

type db struct {
	mu     *sync.RWMutex
	db     *pebble.DB
	pre    []byte
	tbs    map[string]*inline
	closed *atomic.Bool
}

func (d *db) DB() *pebble.DB {
	return d.db
}

func (d *db) Maintain(obj Object) (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isClosed() {
		return ErrClosed
	}
	key := buff.GetBuff()
	val := buff.GetBuff()
	bch := d.db.NewBatch()
	defer func() {
		buff.PutBuff(key)
		buff.PutBuff(val)
		_ = bch.Close()
	}()
	join(key, keySep, d.pre, toBytes(preTBL), toBytes(obj.TableName()))
	oldTbl := &table{Name: obj.TableName(), Index: slices.Clone(obj.TableIndex())}
	newTbl := &table{Name: obj.TableName(), Index: slices.Clone(obj.TableIndex())}
	tin := &inline{Def: newTbl, New: reflectNew(obj)}
	if err = d._getObj(key, oldTbl); err != nil && !errors.Is(err, ErrNotFound) {
		return
	} else if err = d._maintain(bch, tin, oldTbl, newTbl); err != nil {
		return
	} else if err = encode(val, newTbl); err != nil {
		return
	}
	_ = bch.Set(key.Bytes(), val.Bytes(), nil)
	if err = bch.Commit(pebble.Sync); err != nil {
		return err
	} else {
		d.tbs[tin.Def.Name] = tin
		return nil
	}
}

func (d *db) Close() (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isClosed() {
		return ErrClosed
	}
	defer d.closed.Store(true)
	return d.db.Close()
}

func (d *db) Put(obj Object) (id guid.GUID, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isClosed() {
		return id, ErrClosed
	}
	key := buff.GetBuff()
	val := buff.GetBuff()
	idx := buff.GetBuff()
	tin := d.tbs[obj.TableName()]
	tid := obj.TableID()
	bch := d.db.NewBatch()
	defer func() {
		buff.PutBuff(key)
		buff.PutBuff(val)
		buff.PutBuff(idx)
		_ = bch.Close()
	}()
	if tin == nil {
		return id, ErrTableNotFound
	} else if err = encode(val, obj); err != nil {
		return id, err
	}
	d.joinDat(key, tin.Def.Name, tid)
	_ = bch.Delete(key.Bytes(), nil)
	if v, c, e := d.db.Get(key.Bytes()); e != nil && !errors.Is(e, ErrNotFound) {
		return id, e
	} else if e == nil {
		old := buff.GetBuff()
		defer buff.PutBuff(old)
		old.Write(v)
		_ = c.Close()
		o := tin.New()
		if err = decode(old, o); err != nil {
			return id, err
		}
		for _, index := range tin.Def.Index {
			idx.Reset()
			d.joinIdx(idx, tin.Def.Name, index, o.TableField(index), o.TableID())
			_ = bch.Delete(idx.Bytes(), nil)
		}
	}
	_ = bch.Set(key.Bytes(), val.Bytes(), nil)
	for _, index := range tin.Def.Index {
		idx.Reset()
		d.joinIdx(idx, tin.Def.Name, index, obj.TableField(index), tid)
		_ = bch.Set(idx.Bytes(), tid.Bytes(), nil)
	}
	if err = bch.Commit(pebble.Sync); err != nil {
		return id, err
	} else {
		id = tid
		return
	}
}

func (d *db) Get(obj Object, id guid.GUID) (err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.isClosed() {
		return ErrClosed
	}
	key := buff.GetBuff()
	tin := d.tbs[obj.TableName()]
	defer buff.PutBuff(key)
	if tin == nil {
		return ErrTableNotFound
	}
	d.joinDat(key, obj.TableName(), id)
	return d._getObj(key, obj)
}

func (d *db) Del(obj Object, id guid.GUID) (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.isClosed() {
		return ErrClosed
	}
	key := buff.GetBuff()
	val := buff.GetBuff()
	idx := buff.GetBuff()
	tin := d.tbs[obj.TableName()]
	bch := d.db.NewBatch()
	defer func() {
		buff.PutBuff(key)
		buff.PutBuff(val)
		buff.PutBuff(idx)
		_ = bch.Close()
	}()
	if tin == nil {
		return ErrTableNotFound
	}
	d.joinDat(key, obj.TableName(), id)
	o := tin.New()
	if e := d._getObj(key, o); e != nil && !errors.Is(e, ErrNotFound) {
		return e
	} else if errors.Is(e, ErrNotFound) {
		return nil
	}
	_ = bch.Delete(key.Bytes(), nil)
	for _, index := range tin.Def.Index {
		idx.Reset()
		d.joinIdx(idx, obj.TableName(), index, o.TableField(index), id)
		_ = bch.Delete(idx.Bytes(), nil)
	}
	return bch.Commit(pebble.Sync)
}

func (d *db) Has(obj Object, index ...string) (has bool, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.isClosed() {
		return false, ErrClosed
	}
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	tin := d.tbs[obj.TableName()]
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	if tin == nil {
		return has, ErrTableNotFound
	}
	if len(index) == 0 {
		d.joinDatPre(sKey, tin.Def.Name+keySep)
		d.joinDatPre(eKey, tin.Def.Name+keyLmt)
		has, err = d._hasIndex(sKey, eKey)
	} else {
		for _, field := range index {
			sKey.Reset()
			eKey.Reset()
			d.joinIdxValPre(sKey, tin.Def.Name, field, obj.TableField(field))
			d.joinIdxValPre(eKey, tin.Def.Name, field, obj.TableField(field))
			sKey.WriteByte(keySep[0])
			eKey.WriteByte(keyLmt[0])
			if has, err = d._hasIndex(sKey, eKey); err != nil {
				return has, err
			} else if has {
				break
			}
		}
	}
	return
}

func (d *db) Find(obj Object, limit int64, filterMap map[string]Filter) (all []guid.GUID, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.isClosed() {
		return all, ErrClosed
	} else if limit == 0 {
		return all, nil
	}
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	tin := d.tbs[obj.TableName()]
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	if tin == nil {
		return all, ErrTableNotFound
	}
	all = make([]guid.GUID, 0)
	if idxLen := len(filterMap); idxLen > 0 {
		for index, _ := range filterMap {
			if !slices.Contains(tin.Def.Index, index) {
				return all, ErrIndexNotFound
			}
		}
		idxCnt := 0
		for index, filter := range filterMap {
			var idAll []guid.GUID
			var idErr error
			if idxLen > 1 {
				idAll, idErr = d._rangeIdxIDs(tin, -1, index, filter)
			} else {
				idAll, idErr = d._rangeIdxIDs(tin, limit, index, filter)
			}
			if idErr != nil {
				return all, idErr
			} else if idxCnt == 0 {
				all = idAll
			} else {
				all = idIntersect(all, idAll)
			}
			idxCnt++
		}
	} else {
		all, err = d._rangeDatIDs(tin, limit)
	}
	return all, nil
}

// Open db
func Open(dir string, pre []byte) (DB, error) {
	if pdb, err := pebble.Open(dir, nil); err != nil {
		return nil, err
	} else {
		return OpenWithDB(pdb, pre)
	}
}

// OpenWithDB open db with pebble.DB
func OpenWithDB(pdb *pebble.DB, pre []byte) (DB, error) {
	objDB := &db{
		mu:     new(sync.RWMutex),
		db:     pdb,
		pre:    pre,
		tbs:    make(map[string]*inline),
		closed: new(atomic.Bool),
	}
	return objDB, nil
}

/*
  Package method
*/

func (d *db) _rangeIdxIDs(tin *inline, limit int64, index string, filter Filter) (all []guid.GUID, err error) {
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	d.joinIdxPre(sKey, tin.Def.Name, index+keySep)
	d.joinIdxPre(eKey, tin.Def.Name, index+keyLmt)
	preLen := sKey.Len()
	keyLen := 0
	if i, e := d.db.NewIter(&pebble.IterOptions{LowerBound: sKey.Bytes(), UpperBound: eKey.Bytes()}); e != nil {
		return all, e
	} else {
		defer mustClose(i)
		var curr int64
		all = make([]guid.GUID, 0)
		for i.First(); i.Valid() && (limit < 0 || curr < limit); i.Next() {
			key := i.Key()
			keyLen = len(key)
			if value := key[preLen : keyLen-guid.SLen-1]; filter(index, value) {
				id := guid.MustParse(toString(key[keyLen-guid.SLen:]))
				all = append(all, id)
				curr++
			}
		}
	}
	return all, nil
}

func (d *db) _rangeDatIDs(tin *inline, limit int64) (all []guid.GUID, err error) {
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	d.joinDatPre(sKey, tin.Def.Name+keySep)
	d.joinDatPre(eKey, tin.Def.Name+keyLmt)
	preLen := sKey.Len()
	if i, e := d.db.NewIter(&pebble.IterOptions{LowerBound: sKey.Bytes(), UpperBound: eKey.Bytes()}); e != nil {
		return all, e
	} else {
		defer mustClose(i)
		var curr int64
		all = make([]guid.GUID, 0)
		for i.First(); i.Valid() && (limit < 0 || curr < limit); i.Next() {
			key := i.Key()
			curr++
			id := guid.MustParse(toString(key[preLen:]))
			all = append(all, id)
		}
	}
	return all, nil
}

func (d *db) _hasIndex(sk, ek *bytes.Buffer) (has bool, err error) {
	i, e := d.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()})
	if e != nil {
		return has, e
	}
	defer mustClose(i)
	has = i.First()
	return
}

func (d *db) _getObj(key *bytes.Buffer, val any) error {
	buf := buff.GetBuff()
	defer buff.PutBuff(buf)
	if e := d._getVal(buf, key.Bytes()); e != nil {
		return e
	}
	return decode(buf, val)
}

func (d *db) _getVal(buf *bytes.Buffer, key []byte) error {
	if v, c, e := d.db.Get(key); e != nil {
		return e
	} else {
		buf.Write(v)
		return c.Close()
	}
}

func (d *db) _maintain(bch *pebble.Batch, tin *inline, newTbl, oldTbl *table) error {
	sk := buff.GetBuff()
	ek := buff.GetBuff()
	defer func() {
		buff.PutBuff(sk)
		buff.PutBuff(ek)
	}()
	d.joinDatPre(sk, newTbl.Name+keySep)
	d.joinDatPre(ek, newTbl.Name+keyLmt)
	delIdx := make([]string, 0)
	for i := 0; i < len(oldTbl.Index); i++ {
		if !slices.Contains(newTbl.Index, oldTbl.Index[i]) {
			delIdx = append(delIdx, oldTbl.Index[i])
		}
	}
	addIdx := make([]string, 0)
	for i := 0; i < len(newTbl.Index); i++ {
		if !slices.Contains(oldTbl.Index, newTbl.Index[i]) {
			addIdx = append(addIdx, newTbl.Index[i])
		}
	}
	if len(delIdx) == 0 && len(addIdx) == 0 {
		return nil
	}
	var iter *pebble.Iterator
	if i, e := d.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
		return e
	} else {
		iter = i
		defer mustClose(iter)
	}
	tmp := buff.GetBuff()
	defer buff.PutBuff(tmp)
	for iter.First(); iter.Valid(); iter.Next() {
		o := tin.New()
		if b, e := iter.ValueAndErr(); e != nil {
			return e
		} else if err := decode(bytes.NewBuffer(b), o); err != nil {
			return err
		}
		for i := 0; i < len(delIdx); i++ {
			tmp.Reset()
			d.joinIdx(tmp, tin.Def.Name, delIdx[i], o.TableField(delIdx[i]), o.TableID())
			_ = bch.Delete(tmp.Bytes(), nil)
		}
		for i := 0; i < len(addIdx); i++ {
			tmp.Reset()
			d.joinIdx(tmp, tin.Def.Name, addIdx[i], o.TableField(addIdx[i]), o.TableID())
			_ = bch.Set(tmp.Bytes(), o.TableID().Bytes(), nil)
		}
	}
	return nil
}

func (d *db) joinDat(buf *bytes.Buffer, name string, id guid.GUID) {
	// Eg. dat/table/id
	join(buf, keySep, d.pre, toBytes(preDAT), toBytes(name), toBytes(id.String()))
}

func (d *db) joinDatPre(buf *bytes.Buffer, name string) {
	// Eg. dat/table
	join(buf, keySep, d.pre, toBytes(preDAT), toBytes(name))
}

func (d *db) joinIdx(buf *bytes.Buffer, name, index string, value []byte, id guid.GUID) {
	// Eg. idx/table/index/value/id
	join(buf, keySep, d.pre, toBytes(preIDX), toBytes(name), toBytes(index), value, toBytes(id.String()))
}

func (d *db) joinIdxValPre(buf *bytes.Buffer, name, index string, value []byte) {
	// Eg. idx/table/index/value
	join(buf, keySep, d.pre, toBytes(preIDX), toBytes(name), toBytes(index), value)
}

func (d *db) joinIdxPre(buf *bytes.Buffer, name, index string) {
	// Eg. idx/table/index
	join(buf, keySep, d.pre, toBytes(preIDX), toBytes(name), toBytes(index))
}

func (d *db) isClosed() bool {
	return d.closed.Load()
}
