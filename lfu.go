package gocache

import (
	"reflect"
	"sync"
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

func (lfu *Lfu[K, V]) OrderPrint(frequent int) {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	for frequent, lru := range lfu.frequent {
		lru.OrderPrint(frequent)
	}

}

// 为了方便修改， 一样也需要一个双向链表

func (lfu *Lfu[K, V]) add(index int, key K, value V) {
	if _, ok := lfu.frequent[index]; !ok {
		lfu.frequent[index] = &Lru[K, V]{
			lru:  make(map[K]*element[K, V], 0),
			size: lfu.size,
			lock: sync.RWMutex{},
			root: &element[K, V]{},
			last: &element[K, V]{},
		}
	}

	lfu.frequent[index].Add(key, value)
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
	if index, ok := lfu.cache[key]; ok {
		if _, ok := lfu.frequent[index]; ok {
			lfu.frequent[index].Remove(key)
		}
		delete(lfu.cache, key)
	}
}

func (lfu *Lfu[K, V]) Add(key K, value V) (K, bool) {
	// 添加一个key
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	if li, ok := lfu.cache[key]; ok {
		// 判断是否存在新层， 不存在就新建
		lfu.cache[key] = li + 1
		if li/lfu.claddingSize != li+1/lfu.claddingSize {
			// 从原来的那层中删除
			lfu.frequent[li].Remove(key)
			// 原来的那一层没有了就删除
			if lfu.frequent[li].Len() == 0 {
				delete(lfu.frequent, li)
				// 计算最小层
				lfu.getMin(li + 1)
			}
		}
		lfu.add(li+1/lfu.claddingSize, key, value)
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
	if index, ok := lfu.cache[key]; ok {
		if v, ok := lfu.frequent[index]; ok {
			return v.Get(key)
		}
	}
	var v V
	return reflect.Zero(reflect.TypeOf(v)).Interface().(V), false

}
