package gocache

import (
	"testing"
)

func Test_Lru(t *testing.T) {
	l := NewCache[string, int](3, LRU)
	t.Log("------------- apple 1  ---------------")
	l.Add("apple", 1)
	if l.LastKey() != "apple" && l.Len() != 1 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 1 {
		t.Fatal()
	}
	l.OrderPrint()
	t.Log("------------- orange 2 ---------------")
	// time.Sleep(time.Second)
	l.Add("orange", 2)
	if l.LastKey() != "apple" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("orange"); !ok || value != 2 {
		t.Fatal()
	}
	t.Log("------------- apple 3 ---------------")
	l.OrderPrint()
	l.Add("apple", 3)
	if l.LastKey() != "orange" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 3 {
		t.Fatal()
	}
	t.Log("------------- apple 378 ---------------")
	l.OrderPrint()
	l.Add("orange", 378)
	if l.LastKey() != "apple" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 3 {
		t.Fatal()
	}
	t.Log("------------- orange 313 ---------------")
	l.OrderPrint()
	l.Add("orange", 313)
	if l.LastKey() != "apple" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("orange"); !ok || value != 313 {
		t.Fatal()
	}
	t.Log("------------- apple 44 ---------------")
	l.OrderPrint()
	l.Add("apple", 44)
	if l.LastKey() != "orange" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 44 {
		t.Fatal()
	}
	t.Log("------------- mkey 262 ---------------")
	l.OrderPrint()
	l.Add("mkey", 262)
	if l.LastKey() != "orange" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("mkey"); !ok || value != 262 {
		t.Fatal()
	}
	t.Log("------------- banana 262 ---------------")
	l.OrderPrint()
	l.Add("banana", 262)
	if l.LastKey() != "apple" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("banana"); !ok || value != 262 {
		t.Fatal()
	}
	t.Log("------------- apple 3 ---------------")
	l.OrderPrint()
	l.Add("apple", 3)
	if l.LastKey() != "mkey" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 3 {
		t.Fatal()
	}
	t.Log("------------- orange 313 ---------------")
	l.OrderPrint()
	l.Add("orange", 313)
	if l.LastKey() != "banana" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("orange"); !ok || value != 313 {
		t.Fatal()
	}
	t.Log("------------- apple 3 ---------------")
	l.OrderPrint()
	l.Add("apple", 3)
	if l.LastKey() != "banana" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 3 {
		t.Fatal()
	}
}
