package gocache

import "sync"

type Cacher[K comparable, V any] interface {
	Add(key K, value V) (K, bool) // 添加值， 如果返回 k, true 说明有删除值，并返回删除的key
	Remove(key K)                 // 移除k
	Len() int                     // 长度
	OrderPrint(int)               // 顺序打印
	Get(key K) (V, bool)          // 获取值
	LastKey() K                   // 获取最先要删除的key
}

type Algorithm int

const (
	LRU Algorithm = iota
	LFU
	ALFU
	Simple
)

// size: max length of cache
// claddingSize: only use for lfu or alfu, The number of visits and step size are merged into one layer in order to reduce the level, default 1
func NewCache[K comparable, V any](size int, t Algorithm, claddingSize ...int) Cacher[K, V] {
	// 内存足够的话, 可以设置很大, 所有计算都是O(1)
	if size <= 0 {
		size = 2 << 10
	}
	switch t {
	case Simple:
		return &SimpleCache[K, V]{
			cache: make(map[K]V),
			order: make([]K, size),
			mu:    sync.RWMutex{},
			size:  size,
		}
	case LRU:
		return &Lru[K, V]{
			lru:  make(map[K]*element[K, V]),
			size: size,
			lock: sync.RWMutex{},
			root: &element[K, V]{},
			last: &element[K, V]{},
		}
	case LFU:
		cs := 1
		if len(claddingSize) > 0 && claddingSize[0] > 1 {
			cs = 1
		}
		return &Lfu[K, V]{
			frequent: make(map[int]*Lru[K, V]),
			// 这里是根据key来查询在那一层
			cache:        make(map[K]int),
			mu:           sync.RWMutex{},
			size:         size,
			claddingSize: cs,
		}
	case ALFU:
		cs := 1
		if len(claddingSize) > 0 && claddingSize[0] > 1 {
			cs = 1
		}
		alfu := &Alfu[K, V]{
			&Lfu[K, V]{
				frequent: make(map[int]*Lru[K, V]),
				// 这里是根据key来查询在那一层
				cache:        make(map[K]int),
				mu:           sync.RWMutex{},
				size:         size,
				claddingSize: cs,
			},
		}
		go alfu.auto()
		return alfu
	default:
		return nil
	}

}

const DEFAULTCOUNT = 2 << 10
