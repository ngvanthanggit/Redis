package data_structure

import (
	"time"
)

type Obj struct {
	Value interface{}
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64 // epoch time when the key expires
}

func CreateDict() *Dict {
	return &Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
}

func (d *Dict) NewObj(key string, value interface{}, ttlMs int64) *Obj {
	obj := &Obj{
		Value: value,
	}
	if ttlMs > 0 {
		d.SetExpiry(key, uint64(ttlMs))
	} else {
		d.UnSetExpiryIfExists(key)
	}
	return obj
}

func (d *Dict) GetExpiredDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) SetExpiry(key string, ttlMs uint64) {
	d.expiredDictStore[key] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func (d *Dict) UnSetExpiryIfExists(key string) {
	if _, exist := d.expiredDictStore[key]; !exist {
		return
	}
	delete(d.expiredDictStore, key)
}

func (d *Dict) GetExpiry(key string) (uint64, bool) {
	exp, exist := d.expiredDictStore[key]
	return exp, exist
}

func (d *Dict) HasExpired(key string) bool {
	exp, exist := d.GetExpiry(key)
	if !exist {
		return false
	}

	epochNow := uint64(time.Now().UnixMilli())
	if exp <= epochNow {
		return true
	}

	return false
}

func (d *Dict) Set(key string, obj *Obj) {
	d.dictStore[key] = obj
}

func (d *Dict) Del(key string) bool {
	if _, exist := d.dictStore[key]; exist {
		delete(d.dictStore, key)
		delete(d.expiredDictStore, key)
		return true
	}
	return false
}

func (d *Dict) Get(key string) *Obj {
	v := d.dictStore[key]
	if v != nil {
		if d.HasExpired(key) {
			// passive deletion
			d.Del(key)
			return nil
		}
		return v
	}
	return nil
}
