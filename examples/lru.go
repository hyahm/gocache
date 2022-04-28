package main

import "github.com/hyahm/gocache"

func main() {
	c := gocache.NewCache[int, any](100, gocache.LFU)
	c.Add(1, 2)
	c.Add(4, 12)
	c.OrderPrint(0)

}
