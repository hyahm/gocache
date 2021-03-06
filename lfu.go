package gocache

import (
	"fmt"
	"sync"
)

func defaultLfu[K comparable, V any]() *Lfu[K, V] {
	return &Lfu[K, V]{
		layer: make(map[int]*Lru[K, V]),
		// 这里是根据key来查询在那一层
		cache:        make(map[K]int),
		mu:           sync.RWMutex{},
		size:         DEFAULTCOUNT,
		claddingSize: 1,
	}
}

// 以lru为基础

type Lfu[K comparable, V any] struct {
	layer map[int]*Lru[K, V]

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
	for layer := lfu.max; layer >= lfu.min; layer-- {
		if lfu.layer[layer].len == 0 {
			continue
		}
		fmt.Println("layer: ", layer)
		lfu.layer[layer].OrderPrint()

	}
	// lfu.mu.RLock()
	// defer lfu.mu.RUnlock()
	// level := lfu.min
	// for i := 0; i < lfu.Len(); {
	// 	if lfu.layer[level].Len() != 0 {
	// 		fmt.Println("layer: ", level)
	// 		lfu.layer[level].orderPrint()
	// 	}
	// 	i += lfu.layer[level].Len()
	// 	if i < lfu.Len() {
	// 		level = lfu.getNextMin(level + 1)
	// 	}

	// }
}

// 为了方便修改， 一样也需要一个双向链表

func (lfu *Lfu[K, V]) add(level int, key K, value V) {
	if _, ok := lfu.layer[level]; !ok {
		lfu.layer[level] = &Lru[K, V]{
			lru:  make(map[K]*element[K, V], 0),
			size: lfu.size,
			lock: sync.RWMutex{},
			root: &element[K, V]{},
			last: &element[K, V]{},
			len:  0,
		}
	}

	lfu.layer[level].Add(key, value)
}

func (lfu *Lfu[K, V]) getNextMin(start int) int {
	for {
		if _, ok := lfu.layer[start]; ok {
			return start
		}
		start++
	}
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

	if lfu.layer[start].Len() > 0 {
		lfu.min = start
	} else {
		lfu.getMin(start + 1)
	}
}

func (lfu *Lfu[K, V]) Len() int {
	if lfu.layer == nil {
		lfu = defaultLfu[K, V]()
		return 0
	}
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	return len(lfu.cache)
}

// get lastKey
func (lfu *Lfu[K, V]) LastKey() K {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	return lfu.layer[lfu.min].LastKey()
}

func (lfu *Lfu[K, V]) Remove(key K) {
	if lfu.layer == nil {
		lfu = defaultLfu[K, V]()
		return
	}
	lfu.mu.Lock()
	defer lfu.mu.Unlock()
	// 先找到这个key
	if frequent, ok := lfu.cache[key]; ok {
		layer := frequent / lfu.claddingSize
		lfu.layer[layer].Remove(key)
		if layer == lfu.min && lfu.layer[layer].Len() == 0 {
			// 先找大一点的
			lfu.getMin(layer)
		}

		if len(lfu.cache) > 1 {
			delete(lfu.layer, layer)
		}
	}
}

func (lfu *Lfu[K, V]) Resize(size int) {

	if lfu.layer == nil {
		claddingSize := 1
		if lfu.claddingSize != 0 {
			claddingSize = lfu.claddingSize
		}
		lfu = &Lfu[K, V]{
			layer: make(map[int]*Lru[K, V]),
			// 这里是根据key来查询在那一层
			cache:        make(map[K]int),
			mu:           sync.RWMutex{},
			size:         size,
			claddingSize: claddingSize,
		}
		return
	}
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
	if lfu.layer == nil {
		lfu = defaultLfu[K, V]()
		var v K
		return v, false
	}
	// 添加一个key
	lfu.mu.Lock()
	defer lfu.mu.Unlock()
	if frequent, ok := lfu.cache[key]; ok {

		newLayer := (frequent + 1) / lfu.claddingSize
		// 判断是否存在新层， 不存在就新建
		layer := frequent / lfu.claddingSize
		if layer != newLayer {
			// 从原来的那层中删除
			lfu.layer[layer].Remove(key)
			if lfu.layer[layer].Len() == 0 && layer == lfu.min {
				// // 如果这一行没有数据了, 并且是最小的一行 那么计算最小层
				lfu.min = layer + 1
				if len(lfu.cache) > 1 {
					// 至少留一层
					delete(lfu.layer, layer)
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
			removeKey := lfu.layer[lfu.min].removeLast()
			// 删除总缓存
			delete(lfu.cache, removeKey)
			lfu.cache[key] = 0
			lfu.min = 0
			lfu.add(0, key, value)
			return removeKey, true
		}
		lfu.cache[key] = 0
		lfu.min = 0
		lfu.add(0, key, value)
		// 判断是否超过了缓存值
	}
	return key, false
}

//
func (lfu *Lfu[K, V]) Get(key K) (V, bool) {
	if lfu.layer == nil {
		var v V
		return v, false
	}

	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	if frequent, ok := lfu.cache[key]; ok {
		if v, ok := lfu.layer[frequent/lfu.claddingSize]; ok {
			return v.Get(key)
		}
	}
	var v V
	return v, false

}
