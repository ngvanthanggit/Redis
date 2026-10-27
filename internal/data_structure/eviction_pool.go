package data_structure

import (
	"ngvanthanggit/internal/config"
	"sort"
)

type EvictionCandidate struct {
	key            string
	lastAccessTime uint64
}

type EvictionPool struct {
	pool []*EvictionCandidate
}

type ByLastAccessTime []*EvictionCandidate

func (a ByLastAccessTime) Len() int {
	return len(a)
}

func (a ByLastAccessTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i] // this is cool, not like c++
}

func (a ByLastAccessTime) Less(i, j int) bool {
	return a[i].lastAccessTime < a[j].lastAccessTime
}

func NewPool(size int) *EvictionPool {
	return &EvictionPool{
		pool: make([]*EvictionCandidate, size),
	}
}

var ePool *EvictionPool = NewPool(0)

// push a new item into the pool, maintain the item with lastAccessTime in accending order (old items on the left)
// if pool size > EpoolMaxSize, removes the newest item (because it's not likely to be evicted)
func (p *EvictionPool) Push(key string, lastAccessTime uint64) {
	newItem := &EvictionCandidate{
		key:            key,
		lastAccessTime: lastAccessTime,
	}

	exist := false
	for i := 0; i < len(p.pool); i++ {
		if p.pool[i].key == key {
			exist = true
			p.pool[i] = newItem
		}
	}

	if !exist {
		p.pool = append(p.pool, newItem)
	}

	sort.Sort(ByLastAccessTime(p.pool))
	if len(p.pool) > config.EpoolMaxSize {
		lastIndex := len(p.pool) - 1
		p.pool = p.pool[:lastIndex] // remove the last item - newest item
	}
}

func (p *EvictionPool) Pop() *EvictionCandidate {
	if len(p.pool) == 0 {
		return nil
	}

	oldestItem := p.pool[0]
	p.pool = p.pool[1:]
	return oldestItem
}
