package odb

import (
	"bytes"
	"errors"
	"github.com/azeroth-sha/simple/buff"
	"github.com/azeroth-sha/simple/guid"
	"github.com/cockroachdb/pebble"
	"golang.org/x/exp/slices"
	"path"
	"sync"
	"sync/atomic"
)

type store struct {
	mu     *sync.RWMutex
	db     *pebble.DB
	pre    []byte
	tbs    map[string]*inline
	idl    int
	closed *atomic.Bool
}

func (s *store) DB() *pebble.DB {
	return s.db
}

func (s *store) Maintain(obj Object) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return ErrClosed
	}
	key := buff.GetBuff()
	val := buff.GetBuff()
	bch := s.db.NewBatch()
	defer func() {
		buff.PutBuff(key)
		buff.PutBuff(val)
		_ = bch.Close()
	}()
	join(key, keySep, s.pre, toBytes(preTBL), toBytes(obj.TableName()))
	oldTbl := &table{Name: obj.TableName(), Index: slices.Clone(obj.TableIndex())}
	newTbl := &table{Name: obj.TableName(), Index: slices.Clone(obj.TableIndex())}
	tin := &inline{Def: newTbl, New: reflectNew(obj)}
	if err = s._getObj(key, oldTbl); err != nil && !errors.Is(err, ErrNotFound) {
		return
	} else if err = s.maintain(tin, bch, oldTbl, newTbl); err != nil {
		return
	} else if err = encode(val, newTbl); err != nil {
		return
	}
	_ = bch.Set(key.Bytes(), val.Bytes(), nil)
	if err = bch.Commit(pebble.Sync); err != nil {
		return err
	} else {
		s.tbs[tin.Def.Name] = tin
		return nil
	}
}

func (s *store) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return ErrClosed
	}
	defer s.closed.Store(true)
	return s.db.Close()
}

func (s *store) Put(obj Object) (id guid.GUID, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return id, ErrClosed
	}
	tid := obj.TableID()
	tin := s.tbs[obj.TableName()]
	bch := s.db.NewBatch()
	defer mustClose(bch)
	if tin == nil {
		return id, ErrTableNotFound
	}
	if err = s.del(tin, bch, obj, tid); err != nil {
		return id, err
	} else if err = s.set(tin, bch, obj, tid); err != nil {
		return id, err
	} else if err = bch.Commit(pebble.Sync); err != nil {
		return id, err
	} else {
		id = tid
		return
	}
}

func (s *store) Get(obj Object, id guid.GUID) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isClosed() {
		return ErrClosed
	} else if tin := s.tbs[obj.TableName()]; tin == nil {
		return ErrTableNotFound
	} else {
		return s.get(tin, obj, id)
	}
}

func (s *store) Del(obj Object, id guid.GUID) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed() {
		return ErrClosed
	}
	tin := s.tbs[obj.TableName()]
	bch := s.db.NewBatch()
	defer mustClose(bch)
	if tin == nil {
		return ErrTableNotFound
	} else if err = s.del(tin, bch, obj, id); err != nil {
		return err
	} else {
		return bch.Commit(pebble.Sync)
	}
}

func (s *store) Has(obj Object, index ...string) (has bool, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isClosed() {
		return false, ErrClosed
	} else if tin := s.tbs[obj.TableName()]; tin == nil {
		return has, ErrTableNotFound
	} else {
		if len(index) == 0 {
			return s.hasDat(tin)
		}
		return s.hasIdx(tin, obj, index...)
	}
}

func (s *store) Find(obj Object, search *Search) (all []guid.GUID, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isClosed() {
		return all, ErrClosed
	}
	tin := s.tbs[obj.TableName()]
	if tin == nil {
		return all, ErrTableNotFound
	}
	for index, _ := range search.Filter {
		if !slices.Contains(tin.Def.Index, index) {
			return all, ErrIndexNotFound
		}
	}
	all = make([]guid.GUID, 0)
	if len(search.Filter) > 0 {
		for index, _ := range search.Filter {
			if idAll, idErr := s.findIdxIDs(tin, search, index, all); len(idAll) == 0 || idErr != nil {
				return idAll, idErr
			} else {
				all = idAll
			}
		}
	} else {
		all, err = s.findDatIDs(tin, search)
	}
	slices.SortStableFunc(all, func(a, b guid.GUID) bool {
		if search.Desc {
			return a.Gt(b)
		}
		return a.Lt(b)
	})
	if search.Limit > 0 && len(all) > search.Limit {
		all = all[:search.Limit]
	}
	return all, nil
}

// Open db
func Open(dir string, pre []byte) (DB, error) {
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
func OpenWithDB(pdb *pebble.DB, pre []byte) (DB, error) {
	db := &store{
		mu:     new(sync.RWMutex),
		db:     pdb,
		pre:    pre,
		tbs:    make(map[string]*inline),
		idl:    guid.BLen,
		closed: new(atomic.Bool),
	}
	return db, nil
}

/*
  Package method
*/

func (s *store) findIdxIDs(tin *inline, search *Search, index string, inAll []guid.GUID) (all []guid.GUID, err error) {
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	val := buff.GetBuff()
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
		buff.PutBuff(val)
	}()
	s.joinIdxPre(sKey, tin.Def.Name, index+keySep)
	s.joinIdxPre(eKey, tin.Def.Name, index+keyLmt)
	limit := search.Limit
	noLimit := limit == 0 || len(search.Filter) > 1
	noInID := len(inAll) == 0
	filter := search.Filter[index]
	preLen := sKey.Len()
	keyLen := 0
	allLen := 0
	if i, e := s.db.NewIter(&pebble.IterOptions{LowerBound: sKey.Bytes(), UpperBound: eKey.Bytes()}); e != nil {
		return all, e
	} else {
		defer mustClose(i)
		all = make([]guid.GUID, 0)
		var idxKey []byte
		if !search.Desc {
			for i.First(); i.Valid() && (noLimit || allLen < limit); i.Next() {
				bufReset(val)
				idxKey = i.Key()
				if keyLen = len(idxKey); keyLen < preLen+s.idl+1 ||
					idxKey[keyLen-s.idl-1] != keySep[0] {
					continue
				} else if val.Write(idxKey[preLen : keyLen-s.idl-1]); filter(index, val.Bytes()) {
					if id := parseGUID(idxKey[keyLen-s.idl:]); noInID || slices.Contains(inAll, id) {
						all = append(all, id)
						allLen++
					}
				}
			}
		} else {
			for i.Last(); i.Valid() && (noLimit || allLen < limit); i.Prev() {
				bufReset(val)
				idxKey = i.Key()
				if keyLen = len(idxKey); keyLen < preLen+s.idl+1 ||
					idxKey[keyLen-s.idl-1] != keySep[0] {
					continue
				} else if val.Write(idxKey[preLen : keyLen-s.idl-1]); filter(index, val.Bytes()) {
					if id := parseGUID(idxKey[keyLen-s.idl:]); noInID || slices.Contains(inAll, id) {
						all = append(all, id)
						allLen++
					}
				}
			}
		}
	}
	return all, nil
}

func (s *store) findDatIDs(tin *inline, search *Search) (all []guid.GUID, err error) {
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	s.joinDatPre(sKey, tin.Def.Name+keySep)
	s.joinDatPre(eKey, tin.Def.Name+keyLmt)
	limit := search.Limit
	noLimit := limit == 0
	preLen := sKey.Len()
	allLen := 0
	if i, e := s.db.NewIter(&pebble.IterOptions{LowerBound: sKey.Bytes(), UpperBound: eKey.Bytes()}); e != nil {
		return all, e
	} else {
		defer mustClose(i)
		all = make([]guid.GUID, 0)
		if !search.Desc {
			for i.First(); i.Valid() && (noLimit || allLen < limit); i.Next() {
				if key := i.Key(); len(key) != preLen+s.idl {
					continue
				} else {
					allLen++
					all = append(all, parseGUID(key[preLen:]))
				}
			}
		} else {
			for i.Last(); i.Valid() && (noLimit || allLen < limit); i.Prev() {
				if key := i.Key(); len(key) != preLen+s.idl {
					continue
				} else {
					allLen++
					all = append(all, parseGUID(key[preLen:]))
				}
			}
		}
	}
	return all, nil
}

func (s *store) hasIdx(tin *inline, obj Object, list ...string) (has bool, err error) {
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	for _, index := range list {
		bufReset(sKey, eKey)
		s.joinIdxValPre(sKey, tin.Def.Name, index, obj.TableField(index))
		s.joinIdxValPre(eKey, tin.Def.Name, index, obj.TableField(index))
		sKey.WriteByte(keySep[0])
		eKey.WriteByte(keyLmt[0])
		if has, err = s._hasIndex(sKey, eKey); err != nil {
			return has, err
		} else if has {
			break
		}
	}
	return
}

func (s *store) hasDat(tin *inline) (has bool, err error) {
	sKey := buff.GetBuff()
	eKey := buff.GetBuff()
	defer func() {
		buff.PutBuff(sKey)
		buff.PutBuff(eKey)
	}()
	s.joinDatPre(sKey, tin.Def.Name+keySep)
	s.joinDatPre(eKey, tin.Def.Name+keyLmt)
	if i, e := s.db.NewIter(&pebble.IterOptions{LowerBound: sKey.Bytes(), UpperBound: eKey.Bytes()}); e != nil {
		return has, e
	} else {
		defer mustClose(i)
		has = i.First()
	}
	return
}

func (s *store) del(tin *inline, bch *pebble.Batch, _ Object, id guid.GUID) (err error) {
	key := buff.GetBuff()
	defer buff.PutBuff(key)
	s.joinDat(key, tin.Def.Name, id)
	o := tin.New()
	if e := s._getObj(key, o); e != nil && !errors.Is(e, ErrNotFound) {
		return e
	} else if errors.Is(e, ErrNotFound) {
		return
	} else {
		_ = bch.Delete(key.Bytes(), nil)
	}
	for _, index := range tin.Def.Index {
		bufReset(key)
		s.joinIdx(key, tin.Def.Name, index, o.TableField(index), id)
		_ = bch.Delete(key.Bytes(), nil)
	}
	return
}

func (s *store) set(tin *inline, bch *pebble.Batch, obj Object, id guid.GUID) (err error) {
	key := buff.GetBuff()
	val := buff.GetBuff()
	defer func() {
		buff.PutBuff(key)
		buff.PutBuff(val)
	}()
	s.joinDat(key, tin.Def.Name, id)
	if err = encode(val, obj); err != nil {
		return
	} else {
		_ = bch.Set(key.Bytes(), val.Bytes(), nil)
	}
	for _, index := range tin.Def.Index {
		bufReset(key)
		s.joinIdx(key, tin.Def.Name, index, obj.TableField(index), id)
		_ = bch.Set(key.Bytes(), id[:], nil)
	}
	return
}

func (s *store) get(tin *inline, obj Object, id guid.GUID) error {
	key := buff.GetBuff()
	defer buff.PutBuff(key)
	s.joinDat(key, tin.Def.Name, id)
	return s._getObj(key, obj)
}

func (s *store) maintain(tin *inline, bch *pebble.Batch, oldTbl, newTbl *table) error {
	sk := buff.GetBuff()
	ek := buff.GetBuff()
	defer func() {
		buff.PutBuff(sk)
		buff.PutBuff(ek)
	}()
	s.joinDatPre(sk, newTbl.Name+keySep)
	s.joinDatPre(ek, newTbl.Name+keyLmt)
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
	if i, e := s.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
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
			bufReset(tmp)
			s.joinIdx(tmp, tin.Def.Name, delIdx[i], o.TableField(delIdx[i]), o.TableID())
			_ = bch.Delete(tmp.Bytes(), nil)
		}
		for i := 0; i < len(addIdx); i++ {
			bufReset(tmp)
			s.joinIdx(tmp, tin.Def.Name, addIdx[i], o.TableField(addIdx[i]), o.TableID())
			_ = bch.Set(tmp.Bytes(), o.TableID().Bytes(), nil)
		}
	}
	return nil
}

func (s *store) _hasIndex(sk, ek *bytes.Buffer) (has bool, err error) {
	if i, e := s.db.NewIter(&pebble.IterOptions{LowerBound: sk.Bytes(), UpperBound: ek.Bytes()}); e != nil {
		return has, e
	} else {
		defer mustClose(i)
		has = i.First()
	}
	return
}

func (s *store) _getObj(key *bytes.Buffer, val any) error {
	buf := buff.GetBuff()
	defer buff.PutBuff(buf)
	if e := s._getVal(buf, key.Bytes()); e != nil {
		return e
	}
	return decode(buf, val)
}

func (s *store) _getVal(buf *bytes.Buffer, key []byte) error {
	if v, c, e := s.db.Get(key); e != nil {
		return e
	} else {
		buf.Write(v)
		return c.Close()
	}
}

func (s *store) joinDat(buf *bytes.Buffer, name string, id guid.GUID) {
	// Eg. dat/table/id
	join(buf, keySep, s.pre, toBytes(preDAT), toBytes(name), id.Bytes())
}

func (s *store) joinDatPre(buf *bytes.Buffer, name string) {
	// Eg. dat/table
	join(buf, keySep, s.pre, toBytes(preDAT), toBytes(name))
}

func (s *store) joinIdx(buf *bytes.Buffer, name, index string, value []byte, id guid.GUID) {
	// Eg. idx/table/index/value/id
	join(buf, keySep, s.pre, toBytes(preIDX), toBytes(name), toBytes(index), value, id.Bytes())
}

func (s *store) joinIdxValPre(buf *bytes.Buffer, name, index string, value []byte) {
	// Eg. idx/table/index/value
	join(buf, keySep, s.pre, toBytes(preIDX), toBytes(name), toBytes(index), value)
}

func (s *store) joinIdxPre(buf *bytes.Buffer, name, index string) {
	// Eg. idx/table/index
	join(buf, keySep, s.pre, toBytes(preIDX), toBytes(name), toBytes(index))
}

func (s *store) isClosed() bool {
	return s.closed.Load()
}
