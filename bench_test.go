package gocache

import "testing"

type MyValue []byte

func (self MyValue) Size() int {
	return cap(self)
}

func BenchmarkGet(b *testing.B) {
	cache := NewCache[string, []byte](64*1024*1024, LRU)
	value := make(MyValue, 1000)
	cache.Add("stuff", value)
	for i := 0; i < b.N; i++ {
		val, ok := cache.Get("stuff")
		if !ok {
			panic("error")
		}
		_ = val
	}
}
