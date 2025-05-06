package odb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/azeroth-sha/simple/guid"
	"github.com/cockroachdb/pebble"
	"path"
	"slices"
	"sync"
	"sync/atomic"
)

type objectDB struct {
	mu     *sync.RWMutex
	db     *pebble.DB
	pre    []byte
	tblMap map[string]*inline
	closed *atomic.Bool
}

func (o *objectDB) DB() *pebble.DB {
	return o.db
}

func (o *objectDB) Maintain(obj Object) (err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.isClosed() {
		return ErrClosed
	}
	key := getBuf()
	val := getBuf()
	bch := o.getBatch()
	defer func() {
		resetBuf(key, val)
		discardErr(bch.Close)
	}()
	oldTbl := new(table)
	newTbl := &table{Name: obj.TableName(), Index: slices.Clone(obj.TableIndex())}
	tin := &inline{Def: newTbl, New: reflectNew(obj)}
	o.tblKey(key, tin)
	if err = o.getAny(key.Bytes(), oldTbl); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	} else if err = encode(val, newTbl); err != nil {
		return err
	} else if err = o.maintain(tin, bch, oldTbl, newTbl); err != nil {
		return err
	}
	_ = bch.Set(key.Bytes(), val.Bytes(), nil)
	if err = bch.Commit(pebble.Sync); err == nil {
		o.tblMap[tin.Def.Name] = tin
	}
	return err
}

func (o *objectDB) Close() (err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.isClosed() {
		return ErrClosed
	}
	defer o.closed.Store(true)
	return o.db.Close()
}

func (o *objectDB) Put(obj Object) (id guid.GUID, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.isClosed() {
		return id, ErrClosed
	}
	tin := o.tblMap[obj.TableName()]
	if tin == nil {
		return id, ErrTableNotFound
	}
	tid := obj.TableID()
	bch := o.getBatch()
	defer discardErr(bch.Close)
	if err = o.objDel(tin, bch, nil, tid); err != nil && !errors.Is(err, ErrNotFound) {
		return id, err
	} else if err = o.objPut(tin, bch, obj, tid); err != nil {
		return id, err
	} else if err = bch.Commit(pebble.Sync); err != nil {
		return id, err
	} else {
		id = tid
	}
	return id, err
}

func (o *objectDB) Get(obj Object, id guid.GUID) (err error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.isClosed() {
		return ErrClosed
	}
	tin := o.tblMap[obj.TableName()]
	if tin == nil {
		return ErrTableNotFound
	} else {
		return o.objGet(tin, nil, obj, id)
	}
}

func (o *objectDB) Del(obj Object, id guid.GUID) (err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.isClosed() {
		return ErrClosed
	}
	tin := o.tblMap[obj.TableName()]
	if tin == nil {
		return ErrTableNotFound
	}
	bch := o.getBatch()
	defer discardErr(bch.Close)
	if err = o.objDel(tin, bch, obj, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			err = nil
		}
		return err
	} else {
		return bch.Commit(pebble.Sync)
	}
}

func (o *objectDB) Has(obj Object, index ...string) (has bool, err error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.isClosed() {
		return has, ErrClosed
	}
	tin := o.tblMap[obj.TableName()]
	if tin == nil {
		return has, ErrTableNotFound
	} else if len(index) == 0 {
		return o.objHasDat(tin)
	} else if !checkIndex(tin.Def.Index, index...) {
		return has, ErrIndexNotFound
	} else {
		for _, idx := range index {
			if has, err = o.objHasIdx(tin, idx, obj.TableField(idx)); err != nil {
				return has, err
			} else if has {
				return has, nil
			}
		}
	}
	return has, err
}

func (o *objectDB) Find(obj Object, search *Search) (all []guid.GUID, err error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if o.isClosed() {
		return all, ErrClosed
	}
	if search == nil {
		search = new(Search)
	}
	tin := o.tblMap[obj.TableName()]
	if tin == nil {
		return all, ErrTableNotFound
	} else if !checkIndex(tin.Def.Index, search.Index...) {
		return all, ErrIndexNotFound
	}
	all = make([]guid.GUID, 0)
	if len(search.Index) == 0 && search.Filter != nil {
		for _, index := range search.Index {
			if idAll, idErr := o.objFindIDByIndex(tin, search, index, all); len(idAll) == 0 || idErr != nil {
				return all, idErr
			} else {
				all = idAll
			}
		}
	} else {
		if idAll, idErr := o.objFindID(tin, search); len(idAll) == 0 || idErr != nil {
			return all, idErr
		} else {
			all = idAll
		}
	}
	slices.SortStableFunc(all, func(a, b guid.GUID) int {
		if search.Desc {
			return bytes.Compare(b[:], a[:])
		}
		return bytes.Compare(a[:], b[:])
	})
	if search.Limit > 0 && len(all) > search.Limit {
		all = all[:search.Limit]
	}
	return all, err
}

// Open db
func Open(dir string, pre []byte) (ODB, error) {
	opts := &pebble.Options{
		WALDir: path.Join(dir, `./wal`),
	}
	if db, err := pebble.Open(dir, opts); err != nil {
		return nil, err
	} else {
		return OpenWithDB(db, pre)
	}
}

// OpenWithDB open db with pebble.DB
func OpenWithDB(pdb *pebble.DB, pre []byte) (ODB, error) {
	db := &objectDB{
		mu:     new(sync.RWMutex),
		db:     pdb,
		pre:    pre,
		tblMap: make(map[string]*inline),
		closed: new(atomic.Bool),
	}
	return db, nil
}

/*
  Package private functions
*/

func (o *objectDB) objFindIDByIndex(tin *inline, search *Search, index string, inIDs []guid.GUID) (all []guid.GUID, err error) {
	sk := getBuf()
	ek := getBuf()
	defer putBuf(sk, ek)
	if search.UnixL > 0 {
		o.idxKey(sk, tin, index, toGUIDWithSec(search.UnixL, 0x00))
	} else {
		o.idxKey(sk, tin, index, guid.NULL).WriteString(keySep)
	}
	if search.UnixU > 0 {
		o.idxKey(sk, tin, index, toGUIDWithSec(search.UnixU, 0xFF))
	} else {
		o.idxKey(ek, tin, index, guid.NULL).WriteString(keyLmt)
	}
	ignore := len(inIDs) == 0
	var value []byte
	if i, e := o.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
		return all, e
	} else {
		defer discardErr(i.Close)
		all = make([]guid.GUID, 0)
		if !search.Desc {
			for i.First(); i.Valid(); i.Next() {
				key := i.Key()
				id := toGUID(key[len(key)-guid.BLen:])
				if !ignore && !slices.Contains(inIDs, id) {
					continue
				} else if value, err = i.ValueAndErr(); err != nil {
					return all, err
				} else if !search.Filter(id, index, value) {
					continue
				}
				all = append(all, toGUID(key[len(key)-guid.BLen:]))
			}
		} else {
			for i.Last(); i.Valid(); i.Prev() {
				key := i.Key()
				id := toGUID(key[len(key)-guid.BLen:])
				if !ignore && !slices.Contains(inIDs, id) {
					continue
				} else if value, err = i.ValueAndErr(); err != nil {
					return all, err
				} else if !search.Filter(id, index, value) {
					continue
				}
				all = append(all, toGUID(key[len(key)-guid.BLen:]))
			}
		}
	}
	return all, err
}

func (o *objectDB) objFindID(tin *inline, search *Search) (all []guid.GUID, err error) {
	sk := getBuf()
	ek := getBuf()
	defer putBuf(sk, ek)
	if search.UnixL > 0 {
		o.datKey(sk, tin, toGUIDWithSec(search.UnixL, 0x00))
	} else {
		o.datKey(sk, tin, guid.NULL).WriteString(keySep)
	}
	if search.UnixU > 0 {
		o.datKey(sk, tin, toGUIDWithSec(search.UnixU, 0xFF))
	} else {
		o.datKey(ek, tin, guid.NULL).WriteString(keyLmt)
	}
	if i, e := o.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
		return all, e
	} else {
		defer discardErr(i.Close)
		all = make([]guid.GUID, 0)
		if !search.Desc {
			for i.First(); i.Valid(); i.Next() {
				key := i.Key()
				all = append(all, toGUID(key[len(key)-guid.BLen:]))
			}
		} else {
			for i.Last(); i.Valid(); i.Prev() {
				key := i.Key()
				all = append(all, toGUID(key[len(key)-guid.BLen:]))
			}
		}
	}
	return all, err
}

func (o *objectDB) objHasIdx(tin *inline, index string, dst []byte) (bool, error) {
	sk := getBuf()
	ek := getBuf()
	defer putBuf(sk, ek)
	o.idxKey(sk, tin, index, guid.NULL).WriteString(keySep)
	o.idxKey(ek, tin, index, guid.NULL).WriteString(keyLmt)
	if i, e := o.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
		return false, e
	} else {
		defer discardErr(i.Close)
		for i.First(); i.Valid(); i.Next() {
			if bs, er := i.ValueAndErr(); er != nil {
				return false, er
			} else if bytes.Equal(bs, dst) {
				return true, nil
			}
		}
		return false, nil
	}
}

func (o *objectDB) objHasDat(tin *inline) (bool, error) {
	sk := getBuf()
	ek := getBuf()
	defer putBuf(sk, ek)
	o.datKey(sk, tin, guid.NULL).WriteString(keySep)
	o.datKey(ek, tin, guid.NULL).WriteString(keyLmt)
	if i, e := o.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
		return false, e
	} else {
		defer discardErr(i.Close)
		return i.First(), nil
	}
}

func (o *objectDB) objGet(tin *inline, _ *pebble.Batch, obj Object, id guid.GUID) error {
	key := getBuf()
	defer putBuf(key)
	o.datKey(key, tin, id)
	return o.getObj(key.Bytes(), obj)
}

func (o *objectDB) objDel(tin *inline, bch *pebble.Batch, _ Object, id guid.GUID) error {
	key := getBuf()
	defer putBuf(key)
	o.datKey(key, tin, id)
	obj := tin.New()
	if err := o.getObj(key.Bytes(), obj); err != nil {
		return err
	}
	_ = bch.Delete(key.Bytes(), nil)
	for _, idx := range tin.Def.Index {
		resetBuf(key)
		o.idxKey(key, tin, idx, id)
		_ = bch.Delete(key.Bytes(), nil)
	}
	return nil
}

func (o *objectDB) objPut(tin *inline, bch *pebble.Batch, obj Object, id guid.GUID) error {
	key := getBuf()
	val := getBuf()
	defer putBuf(key, val)
	o.datKey(key, tin, id)
	if err := encode(val, obj); err != nil {
		return err
	}
	_ = bch.Set(key.Bytes(), val.Bytes(), nil)
	for _, idx := range tin.Def.Index {
		resetBuf(key)
		o.idxKey(key, tin, idx, obj.TableID())
		_ = bch.Set(key.Bytes(), obj.TableField(idx), nil)
	}
	return nil
}

func (o *objectDB) getObj(key []byte, obj Object) error {
	return o.getAny(key, obj)
}

func (o *objectDB) getAny(key []byte, val any) error {
	buf := getBuf()
	defer putBuf(buf)
	if err := o.dbGet(key, buf); err != nil {
		return err
	}
	return decode(buf, val)
}

func (o *objectDB) dbGet(k []byte, w *bytes.Buffer) error {
	if v, c, e := o.db.Get(k); e != nil {
		return e
	} else {
		defer discardErr(c.Close)
		_, _ = w.Write(v)
		return e
	}
}

func (o *objectDB) maintain(tin *inline, bch *pebble.Batch, oldTbl, newTbl *table) error {
	sk := getBuf()
	ek := getBuf()
	tmp := getBuf()
	defer func() {
		putBuf(sk, ek, tmp)
	}()
	o.datKey(sk, tin, guid.NULL).WriteString(keySep)
	o.datKey(ek, tin, guid.NULL).WriteString(keyLmt)
	delIdx := make([]string, 0)
	addIdx := make([]string, 0)
	for _, idx := range oldTbl.Index {
		if !slices.Contains(newTbl.Index, idx) {
			delIdx = append(delIdx, idx)
		}
	}
	for _, idx := range newTbl.Index {
		if !slices.Contains(oldTbl.Index, idx) {
			addIdx = append(addIdx, idx)
		}
	}
	var iter *pebble.Iterator
	if i, e := o.db.NewIter(&pebble.IterOptions{
		LowerBound: sk.Bytes(),
		UpperBound: ek.Bytes(),
	}); e != nil {
		return e
	} else {
		iter = i
		defer discardErr(iter.Close)
	}
	for iter.First(); iter.Valid(); iter.Next() {
		v := tin.New()
		if bs, er := iter.ValueAndErr(); er != nil {
			return er
		} else if er = decode(bytes.NewBuffer(bs), v); er != nil {
			return er
		} else {
			for _, idx := range delIdx {
				resetBuf(tmp)
				o.idxKey(tmp, tin, idx, v.TableID())
				_ = bch.Delete(tmp.Bytes(), nil)
			}
			for _, idx := range addIdx {
				resetBuf(tmp)
				o.idxKey(tmp, tin, idx, v.TableID())
				_ = bch.Set(tmp.Bytes(), v.TableField(idx), nil)
			}
		}
	}
	return nil
}

func (o *objectDB) idxKey(buf *bytes.Buffer, tin *inline, index string, id guid.GUID) *bytes.Buffer {
	name := tin.Def.Name
	nameLen := fmt.Sprintf(`%02x`, len(name))
	indexLen := fmt.Sprintf(`%02x`, len(index))
	if id.Empty() {
		join(buf, o.pre, keySep, toBts(keyIDX), toBts(nameLen), toBts(name), toBts(indexLen), toBts(index))
	} else {
		join(buf, o.pre, keySep, toBts(keyIDX), toBts(nameLen), toBts(name), toBts(indexLen), toBts(index), id.Bytes())
	}
	return buf
}

func (o *objectDB) datKey(buf *bytes.Buffer, tin *inline, id guid.GUID) *bytes.Buffer {
	name := tin.Def.Name
	nameLen := fmt.Sprintf(`%02x`, len(name))
	if id.Empty() {
		join(buf, o.pre, keySep, toBts(keyDAT), toBts(nameLen), toBts(name))
	} else {
		join(buf, o.pre, keySep, toBts(keyDAT), toBts(nameLen), toBts(name), id.Bytes())
	}
	return buf
}

func (o *objectDB) tblKey(buf *bytes.Buffer, tin *inline) *bytes.Buffer {
	name := tin.Def.Name
	nameLen := fmt.Sprintf(`%02x`, len(name))
	join(buf, o.pre, keySep, toBts(keyTBL), toBts(nameLen), toBts(tin.Def.Name))
	return buf
}

func (o *objectDB) getBatch() *pebble.Batch {
	return o.db.NewBatch()
}

func (o *objectDB) isClosed() bool {
	return o.closed.Load()
}
