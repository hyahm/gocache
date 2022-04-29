package gocache

import (
	"reflect"
	"sync"
	"time"

	"github.com/hyahm/golog"
)

type element[K comparable, V any] struct {
	// 上一个元素和下一个元素
	next, prev *element[K, V]
	// The list to which this element belongs.
	//元素的key
	key K
	// 这个元素的值
	value  V
	update time.Time
}

type Lru[K comparable, V any] struct {
	lru map[K]*element[K, V] //  这里存key 和 元素
	//保存第一个元素
	lock      sync.RWMutex
	root      *element[K, V] // sentinel list element, only &root, root.prev, and root.next are used
	last      *element[K, V] // 最后一个元素
	len       int            // 元素长度
	size      int            // 缓存多少元素
	PrintFunc func(key, value K, update time.Time)
}

func defaultLru[K comparable, V any]() *Lru[K, V] {
	return &Lru[K, V]{
		lru:  make(map[K]*element[K, V], 0),
		size: DEFAULTCOUNT,
		lock: sync.RWMutex{},
		root: &element[K, V]{},
		last: &element[K, V]{},
	}
}

//开始存一个值
func (lru *Lru[K, V]) Add(key K, value V) (K, bool) {
	return lru.add(key, value)
}

// 获取值
func (lru *Lru[K, V]) Get(key K) (V, bool) {
	if lru.lru == nil {
		lru = defaultLru[K, V]()
	}
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	if _, ok := lru.lru[key]; ok {
		return lru.lru[key].value, true
	}
	var v V
	return reflect.Zero(reflect.TypeOf(v)).Interface().(V), false
}

func (lru *Lru[K, V]) GetLastKeyUpdateTime() (K, V, time.Time) {
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	return lru.last.key, lru.last.value, lru.last.update
}

func (lru *Lru[K, V]) NextKey(key K) any {
	if lru.lru == nil {
		lru = defaultLru[K, V]()
	}
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	if value, ok := lru.lru[key]; ok {
		if value.next == nil {
			return nil
		}
		return value.next.key
	}
	return nil
}

func (lru *Lru[K, V]) PrevKey(key K) any {
	if lru.lru == nil {
		lru = defaultLru[K, V]()
	}
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	if value, ok := lru.lru[key]; ok {
		if value.prev == nil {
			return nil
		}
		return value.prev.key
	}
	return nil
}

func (lru *Lru[K, V]) Remove(key K) {
	if lru.lru == nil {
		return
	}
	lru.lock.Lock()
	defer lru.lock.Unlock()
	// 不存在就直接返回
	if _, ok := lru.lru[key]; !ok {
		return
	}
	this := lru.lru[key]
	//如果是第一个元素
	if this == lru.root {
		lru.root = lru.root.next
		lru.root.prev = nil
		// 更新第二个元素的值
		lru.lru[lru.root.key] = lru.root
		delete(lru.lru, key)
		lru.len--
		return
	}
	//如果是最后一个
	if this == lru.last {
		lru.last = lru.last.prev
		lru.last.next = nil
		lru.lru[lru.last.key] = lru.last
		delete(lru.lru, key)
		lru.len--
		return
	}

	// 更改上一个元素的下一个值

	lru.lru[key].prev.next = lru.lru[key].next
	//更新下一个元素的上一个值
	lru.lru[key].next.prev = lru.lru[key].prev
	lru.lru[lru.lru[key].prev.key] = lru.lru[key].prev
	lru.lru[lru.lru[key].next.key] = lru.lru[key].next
	//删除
	delete(lru.lru, key)
	lru.len--
}

func (lru *Lru[K, V]) OrderPrint() {
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	li := lru.root
	for li != nil {
		golog.Infof("key: %v, value: %v, update_time: %s", li.key, li.value, li.update.String())
		li = li.next
	}
}

func (lru *Lru[K, V]) Len() int {
	return lru.len
}

func (lru *Lru[K, V]) Resize(size int) {
	//如果缩小了缓存, 那么可能需要删除后面多余的索引
	if size <= 0 || size == lru.size {
		return
	}
	if size > lru.size {
		lru.size = size
		return
	}
	lru.size = size
	lru.lock.Lock()
	defer lru.lock.Unlock()
	for i := size; i < lru.size; i++ {
		lru.removeLast()
	}
}

// 返回被删除的key, 如果没删除返回nil
func (lru *Lru[K, V]) add(key K, value V) (K, bool) {
	//先要判断是否存在这个key, 存在的话，就将元素移动最开始的位置,
	lru.lock.Lock()
	defer lru.lock.Unlock()
	golog.Error("add key: ", key, ", value: ", value)
	golog.Info("lastkey: ", lru.last.key)
	if _, ok := lru.lru[key]; ok {
		//如果是第一个元素的话, 更新操作
		if lru.root.key == key {
			lru.root.value = value
			lru.root.update = time.Now()

		} else {
			// 否则就插入到开头, 开头的元素后移
			golog.Info("move to prev")
			lru.moveToPrev(key, value)
		}
		return key, false
	} else {
		var removeKey K
		var isremove bool
		if lru.Len() >= lru.size {
			newLastKey := lru.last.prev.key
			removeKey = lru.removeLast()
			lru.lru[newLastKey].next = nil
			lru.last = lru.lru[newLastKey]
			isremove = true
		}
		lru.OrderPrint()
		el := &element[K, V]{
			prev:   nil,
			next:   nil,
			update: time.Now(),
			key:    key,
			value:  value,
		}

		if lru.len == 0 {
			// 更新第一个元素
			lru.root = el
			// 更新最后一个元素
			lru.last = el
			// 更新长度
			lru.len = 1
			// 更新lru
			lru.lru[key] = el
			return key, false
		}
		//如果不存在的话, 直接添加到开头
		// 第二个元素抽取出来， 也就是当前的root
		lru.lru[lru.root.key].prev = el

		el.next = lru.lru[lru.root.key]
		//将开头的元素修改成新的元素
		lru.root = el
		lru.lru[key] = el
		//判断长度是否超过了缓存
		if !isremove {
			lru.len++
		}

		return removeKey, isremove
	}

}

// 移除最后一个, 返回key
func (lru *Lru[K, V]) RemoveLast() K {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	return lru.removeLast()
}

func (lru *Lru[K, V]) removeLast() K {
	lastKey := lru.last.key

	if lru.last.prev != nil {
		lru.last.prev.next = nil
		lru.lru[lru.last.prev.key] = lru.last.prev
	} else {
		// 如果上个元素是空的， 那么说明这是开头的一个元素
		if lru.len == 1 {
			// 如果只有这一个元素， 那么全部清空
			lru.root = nil
			lru.last = nil
			lru.len = 0
			lru.lru = make(map[K]*element[K, V])
			return lastKey
		}
		delete(lru.lru, lastKey)
		lru.len--
	}

	return lastKey
}

func (lru *Lru[K, V]) moveToPrev(key K, value V) {
	// 这里面的元素至少有2个, 否则进不来这里
	// 否则就插入到开头, 开头的元素后移
	//把当前位置元素的上一个元素的下一个元素指向本元素的下一个元素
	//el := &element{}

	if lru.len == 2 {
		//如果是2个元素
		//也就是更换元素的值就好了
		//把第一个元素换到第二去
		// 新的root 和 last
		rootKey := lru.last.key
		lastKey := lru.root.key

		lru.lru[lastKey].next = nil
		lru.lru[lastKey].prev = lru.lru[rootKey]

		lru.lru[rootKey].next = lru.lru[lastKey]
		lru.lru[rootKey].prev = nil
		lru.lru[rootKey].value = value
		lru.lru[rootKey].update = time.Now()

		lru.root = lru.lru[rootKey]
		lru.last = lru.lru[lastKey]
		return
	}
	if lru.len > 2 {
		// 拿到这个key的值
		secondKey := lru.root.key
		if key == lru.last.key {
			lru.lru[lru.last.prev.key].next = nil
			lru.last = lru.lru[lru.last.prev.key]
			// 最后一个元素 是最后一个元素
		} else {
			// 如果不是最后一个也不是最前面一个， 一定在中间， 那么直接把
			// 上一个key 的next 等于下一个
			// 下一个的prev 等于上一个
			lru.lru[lru.lru[key].prev.key].next = lru.lru[lru.lru[key].next.key]
			lru.lru[lru.lru[key].next.key].prev = lru.lru[lru.lru[key].prev.key]
		}
		// 移动到开头
		lru.lru[secondKey].prev = lru.lru[key]
		lru.lru[key].next = lru.lru[secondKey]
		lru.lru[key].value = value
		lru.root = lru.lru[key]
	}
}

func (lru *Lru[K, V]) FirstKey() any {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	return lru.root.key
}

func (lru *Lru[K, V]) LastKey() K {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	return lru.last.key
}

func (lru *Lru[K, V]) Clean(n int) {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	lru = nil
	lru.lru = nil
	lru = &Lru[K, V]{
		lru:  make(map[K]*element[K, V], 0),
		len:  0,
		size: n,
		lock: sync.RWMutex{},
		root: &element[K, V]{},
		last: &element[K, V]{},
	}
}

func (lru *Lru[K, V]) Exsit(key K) bool {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if _, ok := lru.lru[key]; ok {
		return true
	}
	return false
}
