package gocache

import (
	"fmt"
	"reflect"
	"sync"
)

// 以lru为基础

type Lfu[K comparable, V any] struct {
	row map[int]*Lru[K, V]

	cache        map[K]int // 保存key得访问量
	min          int       // 记录当前最小层的值
	max          int       // 记录当前最大层的值
	mu           sync.RWMutex
	size         int // 大小
	claddingSize int
	layerOrder   []int // 层从大到小排序
}

func (lfu *Lfu[K, V]) OrderPrint() {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	for layer := lfu.max; layer <= lfu.min; layer-- {
		fmt.Println("layer: ", layer)
		lfu.row[layer].OrderPrint()
	}

}

// 为了方便修改， 一样也需要一个双向链表

func (lfu *Lfu[K, V]) add(level int, key K, value V) {
	if _, ok := lfu.row[level]; !ok {
		lfu.row[level] = &Lru[K, V]{
			lru:  make(map[K]*element[K, V], 0),
			size: lfu.size,
			lock: sync.RWMutex{},
			root: &element[K, V]{},
			last: &element[K, V]{},
			len:  0,
		}
	}

	lfu.row[level].Add(key, value)
}

func (lfu *Lfu[K, V]) getMin(start int) {
	if len(lfu.cache) == 1 {
		lfu.min = start
		return
	}

	if start == lfu.max {
		lfu.min = lfu.max
		return
	}

	if lfu.row[start].Len() > 0 {
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
	return lfu.row[lfu.min].LastKey()
}

func (lfu *Lfu[K, V]) Remove(key K) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()
	// 先找到这个key
	if frequent, ok := lfu.cache[key]; ok {
		level := frequent / lfu.claddingSize
		lfu.row[level].Remove(key)
		if lfu.row[level].Len() == 0 {
			// 先找大一点的
			lfu.getMin(level)
			if len(lfu.cache) > 1 {
				delete(lfu.row, level)
			}
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
		newLayer := (frequent + 1) / lfu.claddingSize
		// 判断是否存在新层， 不存在就新建
		layer := frequent / lfu.claddingSize
		if layer != newLayer {
			// 从原来的那层中删除
			lfu.row[layer].Remove(key)

			if lfu.row[layer].Len() == 0 && layer == lfu.min {
				// // 如果这一行没有数据了, 并且是最小的一行 那么计算最小层
				lfu.getMin(layer)
				if len(lfu.cache) > 1 {
					// 至少留一层
					delete(lfu.row, layer)
				}
			}
		}
		lfu.add(newLayer, key, value)
		if newLayer > lfu.max {
			lfu.max = newLayer
		}
		lfu.cache[key] = frequent + 1
	} else {
		// 如果当前的大小大于等于
		if len(lfu.cache) >= int(lfu.size) {
			// 删除最后一个
			removeKey := lfu.row[lfu.min].RemoveLast()
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
		if v, ok := lfu.row[frequent/lfu.claddingSize]; ok {
			return v.Get(key)
		}
	}
	var v V
	return reflect.Zero(reflect.TypeOf(v)).Interface().(V), false

}
