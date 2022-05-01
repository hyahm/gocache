package gocache

import (
	"fmt"
	"reflect"
	"sync"
)

// 以lru为基础

type Lfu[K comparable, V any] struct {
	row map[int]*Lru[K, V]

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
	level := lfu.min
	fmt.Println("length: ", lfu.Len())
	for i := 0; i < lfu.Len(); {

		if lfu.row[level].Len() != 0 {
			fmt.Println("row: ", level)
			lfu.row[level].orderPrint()
		}
		i += lfu.row[level].Len()
		if i < lfu.Len() {
			level = lfu.getNextMin(level + 1)
		}

	}
	fmt.Println("-------------------------")
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

func (lfu *Lfu[K, V]) getNextMin(start int) int {
	for {
		if _, ok := lfu.row[start]; ok {
			return start
		}
		start++
	}
}

func (lfu *Lfu[K, V]) getMin(start int) {
	fmt.Println("start: ", start)
	if len(lfu.cache) == 1 {
		lfu.min = start
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
		if level == lfu.min && lfu.row[level].Len() == 0 {
			// 先找大一点的
			lfu.getMin(level)
		}

		if len(lfu.cache) > 1 {
			delete(lfu.row, level)
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
		level := frequent / lfu.claddingSize
		fmt.Println("key: ", key)
		fmt.Println("frequent: ", frequent)
		if level != frequent+1/lfu.claddingSize {
			// 从原来的那层中删除
			lfu.row[level].Remove(key)
			fmt.Println("remove key and min: ", lfu.min)
			fmt.Println("remove key and level: ", level)
			if lfu.row[level].Len() == 0 && level == lfu.min {
				// // 如果这一行没有数据了, 并且是最小的一行 那么计算最小层
				lfu.min = level + 1
				if len(lfu.cache) > 1 {
					// 至少留一层
					delete(lfu.row, level)
				}

			}
		}
		lfu.add(frequent+1/lfu.claddingSize, key, value)
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

	fmt.Println("min: ", lfu.min)
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
