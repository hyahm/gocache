package gocache

import (
	"log"
	"reflect"
	"sync"
)

// 简单缓存， 先进先出

type SimpleCache[K comparable, V any] struct {
	cache map[K]V
	order []K
	mu    sync.RWMutex
	size  int
}

func (sc *SimpleCache[K, V]) Add(key K, value V) (K, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if _, ok := sc.cache[key]; ok {
		sc.cache[key] = value
		return key, false
	} else {
		sc.cache[key] = value
		sc.order = append(sc.order, key)
		if len(sc.order) >= sc.size {
			rmKey := sc.order[sc.size-1]
			sc.order = sc.order[1:]
			return rmKey, true
		}
		return key, false
	}
}

func (sc *SimpleCache[K, V]) Get(key K) (V, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if v, ok := sc.cache[key]; ok {
		return v, ok
	}
	var v V

	return reflect.Zero(reflect.TypeOf(v)).Interface().(V), false
}

func (sc *SimpleCache[K, V]) Resize(size int) {
	if size <= 0 || size == sc.size {
		return
	}
	if size > sc.size {
		sc.size = size
		return
	}
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	sc.size = size
	for index, v := range sc.order[size:] {
		sc.order = append(sc.order[:index], sc.order[index+1:]...)
		delete(sc.cache, v)
		return
	}
}

func (sc *SimpleCache[K, V]) Len() int {
	return sc.size
}

func (sc *SimpleCache[K, V]) OrderPrint() {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	for _, key := range sc.order {
		log.Printf("key: %v, value: %v\n", key, sc.cache[key])
	}
}

func (sc *SimpleCache[K, V]) LastKey() K {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.order[sc.size-1]
}

func (sc *SimpleCache[K, V]) Remove(key K) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	for index, v := range sc.order {
		sc.order = append(sc.order[:index], sc.order[index+1:]...)
		delete(sc.cache, v)
		return
	}

}
