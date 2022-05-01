package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	cache := gocache.NewCache[string, string](3, gocache.LFU)
	cache.Add("dddd", "bbbbb")
	// cache.OrderPrint()
	// 	layer:  0
	// key: dddd, value: bbbbb, update_time: ...
	cache.Add("cccc", "111111")
	// cache.OrderPrint()
	// 	layer:  0
	// key: cccc, value: 111111, update_time: ...
	// key: dddd, value: bbbbb, update_time: ...
	cache.Add("cccc", "111111")
	// cache.OrderPrint()
	// 	layer:  0
	// key: dddd, value: bbbbb, update_time: ...
	// 	layer:  1
	// key: cccc, value: 111111, update_time: ...
	cache.Add("aaaa", "111111")
	// cache.OrderPrint()
	// 	layer:  0
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:40:35.380444 +0800 CST m=+0.001039201
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:40:35.380444 +0800 CST m=+0.001039201
	// layer:  1
	// key: cccc, value: 111111, update_time: 2022-05-01 20:40:35.380444 +0800 CST m=+0.001039201
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:41:12.3524927 +0800 CST m=+0.000518001
	// layer:  1
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:41:12.3524927 +0800 CST m=+0.000518001
	// key: cccc, value: 111111, update_time: 2022-05-01 20:41:12.3524927 +0800 CST m=+0.000518001
	cache.Add("aaaa", "111111")
	// cache.OrderPrint()
	// layer:  0
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:41:32.8196747 +0800 CST m=+0.001045201
	// layer:  1
	// key: cccc, value: 111111, update_time: 2022-05-01 20:41:32.8196747 +0800 CST m=+0.001045201
	// layer:  2
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:41:32.8196747 +0800 CST m=+0.001045201
	cache.Add("aaaa", "111111")
	// cache.OrderPrint()
	// layer:  0
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:41:47.8669951 +0800 CST m=+0.001038101
	// layer:  1
	// key: cccc, value: 111111, update_time: 2022-05-01 20:41:47.8669951 +0800 CST m=+0.001038101
	// layer:  3
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:41:47.8669951 +0800 CST m=+0.00103810
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:42:01.3710071 +0800 CST m=+0.001041701
	// layer:  1
	// key: cccc, value: 111111, update_time: 2022-05-01 20:42:01.3710071 +0800 CST m=+0.001041701
	// layer:  4
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:42:01.3710071 +0800 CST m=+0.001041701

	cache.Add("bbb", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: bbb, value: 111111, update_time: 2022-05-01 20:42:24.0629599 +0800 CST m=+0.001001301
	// layer:  1
	// key: cccc, value: 111111, update_time: 2022-05-01 20:42:24.0629599 +0800 CST m=+0.001001301
	// layer:  4
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:42:24.0629599 +0800 CST m=+0.001001301
	cache.Add("ddd", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: ddd, value: 111111, update_time: 2022-05-01 20:42:43.8037083 +0800 CST m=+0.001041501
	// layer:  1
	// key: cccc, value: 111111, update_time: 2022-05-01 20:42:43.8037083 +0800 CST m=+0.001041501
	// layer:  4
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:42:43.8037083 +0800 CST m=+0.001041501

	if value, ok := cache.Get("cccc"); ok {
		fmt.Println(value) // "111111"
	}

}
