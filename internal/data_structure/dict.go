package data_structure

import (
	"log"
	"ngvanthanggit/internal/config"
	"time"
)

type Obj struct {
	Value          interface{}
	LassAccessTime uint32
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64 // epoch time when the key expires
}

func now() uint32 {
	return uint32(time.Now().Unix())
}

func CreateDict() *Dict {
	return &Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
}

func (d *Dict) NewObj(key string, value interface{}, ttlMs int64) *Obj {
	obj := &Obj{
		Value:          value,
		LassAccessTime: now(),
	}
	if ttlMs > 0 {
		d.SetExpiry(key, uint64(ttlMs))
	} else {
		d.UnSetExpiryIfExists(key)
	}
	return obj
}

func (d *Dict) GetDictStore() map[string]*Obj {
	return d.dictStore
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

func (d *Dict) evictRandom() {
	log.Println("trigger randome eviction")
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	for k := range d.dictStore {
		d.Del(k)
		evictCount--
		if evictCount == 0 {
			break
		}
	}
}

func (d *Dict) evict() {
	switch config.EvictionPolicy {
	case "allkeys-random":
		d.evictRandom()
	case "allkeys-lru":
		d.evictLru()
	}
}

func (d *Dict) Set(key string, obj *Obj) {
	d.dictStore[key] = obj
}

func (d *Dict) Del(key string) bool {
	log.Printf("Delete key %s", key)
	if _, exist := d.dictStore[key]; exist {
		delete(d.dictStore, key)
		delete(d.expiredDictStore, key)
		HashKeySpaceStat.Key--
		return true
	}
	return false
}

func (d *Dict) Get(key string) *Obj {
	if len(d.dictStore) == config.MaxKeyNumber {
		d.evict()
	}
	v := d.dictStore[key]
	if v == nil {
		HashKeySpaceStat.Key++
	} else {
		v.LassAccessTime = now()
		if d.HasExpired(key) {
			// passive deletion
			// d.Del(key)
			return nil
		}
		return v
	}
	return nil
}

func (d *Dict) populateEpool() {
	remain := config.EpoolLruSampleSize
	for k := range d.dictStore {
		ePool.Push(k, uint64(d.dictStore[k].LassAccessTime))
		remain--
		if remain <= 0 {
			break
		}
	}
	log.Println("EPool:")
	for _, item := range ePool.pool {
		log.Println(item.key, item.lastAccessTime)
	}
}

func (d *Dict) evictLru() {
	d.populateEpool()
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("trigger LRU eviction")

	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item != nil {
			d.Del(item.key)
		}
	}
}
