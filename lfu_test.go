package gocache

import (
	"testing"
)

func Test_Lfu(t *testing.T) {
	l := NewCache[string, int](3, LFU)
	l.Add("apple", 1)
	if l.LastKey() != "apple" && l.Len() != 1 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 1 {
		t.Fatal()
	}
	// time.Sleep(time.Second)

	l.OrderPrint()
	t.Log("orange")
	l.Add("orange", 2)
	if l.LastKey() != "apple" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("orange"); !ok || value != 2 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("apple", 3)
	if l.LastKey() != "orange" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 3 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("apple", 378)
	if l.LastKey() != "orange" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 378 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("apple", 313)
	if l.LastKey() != "orange" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 313 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("orange", 44)
	if l.LastKey() != "orange" && l.Len() != 2 {
		t.Fatal()
	}
	if value, ok := l.Get("orange"); !ok || value != 44 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("mkey", 262)
	if l.LastKey() != "mkey" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("mkey"); !ok || value != 262 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("banana", 262)
	if l.LastKey() != "banana" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("banana"); !ok || value != 262 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("orange", 313)
	if l.LastKey() != "banana" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("orange"); !ok || value != 313 {
		t.Fatal()
	}

	l.OrderPrint()
	l.Add("apple", 3)
	if l.LastKey() != "banana" && l.Len() != 3 {
		t.Fatal()
	}
	if value, ok := l.Get("apple"); !ok || value != 3 {
		t.Fatal()
	}
}
