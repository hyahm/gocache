package gocache

import (
	"reflect"
	"sync"

	"github.com/hyahm/golog"
)

// 以lru为基础

type Lfu[K comparable, V any] struct {
	frequent map[int]*Lru[K, V]

	// 这里是根据key来查询在那一层
	cache        map[K]int
	min          int // 记录当前最小层的值
	mu           sync.RWMutex
	size         int // 大小
	claddingSize int
}

func (lfu *Lfu[K, V]) OrderPrint() {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	for frequent, lru := range lfu.frequent {
		golog.Info("frequent: ", frequent)
		lru.OrderPrint()
	}

}

// 为了方便修改， 一样也需要一个双向链表

func (lfu *Lfu[K, V]) add(level int, key K, value V) {
	if _, ok := lfu.frequent[level]; !ok {
		lfu.frequent[level] = &Lru[K, V]{
			lru:  make(map[K]*element[K, V], 0),
			size: lfu.size,
			lock: sync.RWMutex{},
			root: &element[K, V]{},
			last: &element[K, V]{},
		}
	}

	lfu.frequent[level].Add(key, value)
}

func (lfu *Lfu[K, V]) getMin(start int) {
	if lfu.frequent[start].Len() > 0 {
		lfu.min = start
	} else {
		lfu.getMin(start + 1)
	}
}

func (lfu *Lfu[K, V]) Len() int {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	return len(lfu.cache)
}

// get lastKey
func (lfu *Lfu[K, V]) LastKey() K {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	return lfu.frequent[lfu.min].LastKey()
}

func (lfu *Lfu[K, V]) Remove(key K) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()
	// 先找到这个key
	if frequent, ok := lfu.cache[key]; ok {
		level := frequent / lfu.claddingSize
		lfu.frequent[level].Remove(key)
		delete(lfu.cache, key)
		if lfu.frequent[level].Len() == 0 {
			delete(lfu.frequent, level)
			lfu.getMin(level + 1)
		}
	}
}

func (lfu *Lfu[K, V]) Resize(size int) {
	if size <= 0 || size == lfu.size {
		return
	}
	if size > lfu.size {
		lfu.size = size
		return
	}
	lfu.size = size
	for i := size; i < lfu.size; i++ {
		rmKey := lfu.LastKey()
		lfu.Remove(rmKey)
	}
}

func (lfu *Lfu[K, V]) Add(key K, value V) (K, bool) {
	// 添加一个key
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	if frequent, ok := lfu.cache[key]; ok {
		// 判断是否存在新层， 不存在就新建
		lfu.cache[key] = frequent + 1
		if frequent/lfu.claddingSize != frequent+1/lfu.claddingSize {
			// 从原来的那层中删除
			lfu.frequent[frequent].Remove(key)
			// 原来的那一层没有了就删除
			if lfu.frequent[frequent].Len() == 0 {
				delete(lfu.frequent, frequent)
				// 计算最小层
				lfu.getMin(frequent + 1)
			}
		}
		lfu.add(frequent+1/lfu.claddingSize, key, value)
	} else {
		// 如果当前的大小大于等于
		if len(lfu.cache) >= int(lfu.size) {
			// 删除最后一个
			removeKey := lfu.frequent[lfu.min].RemoveLast()
			// 删除总缓存
			delete(lfu.cache, removeKey)
		}
		lfu.cache[key] = 0
		lfu.min = 0
		lfu.add(lfu.cache[key]/lfu.claddingSize, key, value)
		// 判断是否超过了缓存值
	}
	return key, false
}

//
func (lfu *Lfu[K, V]) Get(key K) (V, bool) {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	if frequent, ok := lfu.cache[key]; ok {
		if v, ok := lfu.frequent[frequent/lfu.claddingSize]; ok {
			return v.Get(key)
		}
	}
	var v V
	return reflect.Zero(reflect.TypeOf(v)).Interface().(V), false

}
