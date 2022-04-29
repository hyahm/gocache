package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	cache := gocache.NewCache[string, string](3, gocache.LFU)

	cache.Add("dddd", "bbbbb")

	// 1    dddd
	cache.Add("cccc", "111111")

	// 1    cccc  dddd
	cache.Add("cccc", "111111")

	// 1    dddd
	// 2    cccc
	cache.Add("aaaa", "111111")

	// 1  aaaa  dddd
	// 2    cccc
	cache.Add("aaaa", "111111")

	// 1   dddd
	// 2   aaaa cccc
	cache.Add("aaaa", "111111")

	// 1   dddd
	// 2   cccc
	// 3   aaaa

	cache.Add("aaaa", "111111")
	cache.OrderPrint()
	cache.Add("aaaa", "111111")
	cache.Add("aaaa", "111111")
	// 1   dddd
	// 2   cccc
	// 6   aaaa
	cache.Add("bbb", "111111")
	// 1   bbb
	// 2   cccc
	// 6   aaaa
	cache.Add("ddd", "111111")
	// 1   ddd
	// 2   cccc
	// 6   aaaa

	/*
		key:  cccc, value: 111111
		key:  adsf, value: bbbbb
	*/

	if value, ok := cache.Get("cccc"); ok {
		fmt.Println(value) // "111111"
	}

}
