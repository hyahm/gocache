package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	cache := gocache.NewCache[string, string](3, gocache.LFU, 2)

	cache.Add("adsf", "bbbbb")
	cache.Add("cccc", "111111")
	cache.OrderPrint()
	/*
		key:  cccc, value: 111111
		key:  adsf, value: bbbbb
	*/
	if value, ok := cache.Get("cccc"); ok {
		fmt.Println(value) // "111111"
	}

}
