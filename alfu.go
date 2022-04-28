package gocache

import (
	"time"
)

// 基于 lfu 算法， 自动计算 访问频率

type Alfu[K comparable, V any] struct {
	*Lfu[K, V]
}

// 唯一不同的多了一个自动计算 频率的goroutine
func (lfu *Alfu[K, V]) auto() {

	tick := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-tick.C:
			lfu.mu.Lock()
			defer lfu.mu.Unlock()
			for index, lru := range lfu.frequent {
				if index == lfu.min {
					continue
				}
				key, value, update_time := lru.GetLastKeyUpdateTime()
				if time.Since(update_time).Hours() >= 24 {
					// 如果最后访问的时间大于1天时间了， 那么将访问频率减少一半
					// 删除
					lru.Remove(key)
					// 添加到新层中
					lfu.cache[key] = index / 2
					newLevel := lfu.cache[key] / lfu.claddingSize
					if lfu.min > newLevel {
						lfu.min = newLevel
					}

					lfu.add(newLevel, key, value)

				}
			}
		}

	}

}

// func (lfu *Alfu[K, V]) OrderPrint(frequent int) {
// 	lfu.mu.RLock()
// 	defer lfu.mu.RUnlock()
// 	for frequent, lru := range lfu.frequent {
// 		fmt.Printf("%#v\n", lru)
// 		lru.OrderPrint(frequent)
// 	}

// }

// // 为了方便修改， 一样也需要一个双向链表

// func (lfu *Alfu[K, V]) add(index int, key K, value V) {
// 	if _, ok := lfu.frequent[index]; !ok {
// 		lfu.frequent[index] = &Lru[K, V]{
// 			lru:  make(map[K]*element[K, V], 0),
// 			size: lfu.size,
// 			lock: sync.RWMutex{},
// 			root: &element[K, V]{},
// 			last: &element[K, V]{},
// 		}
// 	}

// 	lfu.frequent[index].Add(key, value)
// }

// func (lfu *Alfu[K, V]) getMin(start int) int {
// 	if lfu.frequent[start].Len() > 0 {
// 		return start
// 	} else {
// 		return lfu.getMin(start + 1)
// 	}
// }

// func (lfu *Alfu[K, V]) Len() int {
// 	lfu.mu.RLock()
// 	defer lfu.mu.RUnlock()
// 	return len(lfu.cache)
// }

// // get lastKey
// func (lfu *Alfu[K, V]) LastKey() K {
// 	lfu.mu.RLock()
// 	defer lfu.mu.RUnlock()
// 	return lfu.frequent[lfu.min].LastKey()
// }

// func (lfu *Alfu[K, V]) Remove(key K) {
// 	lfu.mu.Lock()
// 	defer lfu.mu.Unlock()
// 	// 先找到这个key
// 	if index, ok := lfu.cache[key]; ok {
// 		if _, ok := lfu.frequent[index]; ok {
// 			lfu.frequent[index].Remove(key)
// 		}
// 		delete(lfu.cache, key)
// 	}
// }

// func (lfu *Alfu[K, V]) Add(key K, value V) (K, bool) {
// 	// 添加一个key
// 	lfu.mu.Lock()
// 	defer lfu.mu.Unlock()

// 	if li, ok := lfu.cache[key]; ok {
// 		// 如果存在的话，删除此层的值，
// 		lfu.frequent[li].Remove(key)
// 		//添加到新层中
// 		// 判断是否存在新层， 不存在就新建
// 		lfu.cache[key] = li + 1
// 		lfu.add(li+1, key, value)
// 	} else {
// 		lfu.cache[key] = 1
// 		lfu.min = 1
// 		lfu.add(1, key, value)
// 		// 判断是否超过了缓存值
// 		if len(lfu.cache) >= lfu.size {
// 			// 删除最后一个
// 			removeKey := lfu.frequent[lfu.min].RemoveLast()
// 			// 删除总缓存
// 			delete(lfu.cache, removeKey)
// 			if lfu.frequent[lfu.min].Len() == 0 {
// 				// 如果长度为空， 我们就要重新获取最小层
// 				// delete(frequent, min)
// 				// 继续取最小层数
// 				lfu.min = lfu.getMin(lfu.min + 1)
// 			}
// 			return removeKey, true
// 		}
// 	}
// 	return key, false
// }

// //
// func (lfu *Alfu[K, V]) Get(key K) (V, bool) {
// 	lfu.mu.RLock()
// 	defer lfu.mu.RUnlock()
// 	if index, ok := lfu.cache[key]; ok {
// 		if v, ok := lfu.frequent[index]; ok {
// 			return v.Get(key)
// 		}
// 	}
// 	var v V
// 	return reflect.Zero(reflect.TypeOf(v)).Interface().(V), false
// }
