package main

import (
	"github.com/hyahm/gocache"
	"github.com/hyahm/golog"
)

func main() {
	defer golog.Sync()
	l := gocache.NewCache[string, any](3, gocache.LRU)
	l.Add("apple", 1)
	l.OrderPrint()
	// time.Sleep(time.Second)
	l.Add("orange", 2)
	l.OrderPrint()
	l.Add("apple", 3)
	l.OrderPrint()
	l.Add("orange", 378)
	l.OrderPrint()
	l.Add("orange", 313)
	l.OrderPrint()
	l.Add("apple", 262)
	l.OrderPrint()
	l.Add("mkey", 262)
	l.OrderPrint()
	l.Add("banana", 262)
	l.OrderPrint()
	l.Add("apple", 3)
	l.OrderPrint()
	l.Add("orange", 313)
	l.OrderPrint()

}
