package gocache

import (
	"testing"
)

func Test_Add(t *testing.T) {
	l := NewCache[string, int](3, LRU)
	l.Add("apple", 1)
	// time.Sleep(time.Second)
	l.Add("orange", 2)
	l.Add("apple", 3)
	l.Add("orange", 378)
	l.Add("orange", 313)
	l.Add("apple", 262)
	l.Add("mkey", 262)
	l.Add("banana", 262)
	t.Log(l.Len())
	l.Add("apple", 3)
	t.Log(l.Len())
	l.Add("orange", 313)
	t.Log(l.Len())
	l.OrderPrint()
	// x := l.Keys()
	// golog.Info(x)
	// golog.Info(l.Get("orange"))
}
