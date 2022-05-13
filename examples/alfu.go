package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	// only ALFU have Reduce function
	cache := gocache.NewCache[string, string](3, gocache.ALFU).(*gocache.Alfu[string, string])
	cache.Add("dddd", "bbbbb")
	cache.OrderPrint()
	fmt.Println("-----------11111---------------")
	// 	layer:  0
	// key: dddd, value: bbbbb, update_time: ...
	cache.Add("cccc", "111111")

	cache.OrderPrint()
	// 	layer:  0
	// key: cccc, value: 111111, update_time: ...
	// key: dddd, value: bbbbb, update_time: ...
	fmt.Println("-------------222222-------------")
	cache.Add("cccc", "111111")
	cache.OrderPrint()
	fmt.Println("-----------333333---------------")
	// 	layer:  0
	// key: dddd, value: bbbbb, update_time: ...
	// 	layer:  1
	// key: cccc, value: 111111, update_time: ...
	cache.Reduce()
	cache.OrderPrint()
	fmt.Println("-------------44444-------------")
	// 	layer:  0
	// key: cccc, value: 111111, update_time: 2022-05-01 20:54:43.9150871 +0800 CST m=+0.001028301
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:54:43.9150871 +0800 CST m=+0.001028301
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:03:36.1240546 +0800 CST m=+0.001039601
	// key: cccc, value: 111111, update_time: 2022-05-01 21:03:36.1240546 +0800 CST m=+0.001039601
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:03:36.1240317 +0800 CST m=+0.001016701
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: cccc, value: 111111, update_time: 2022-05-01 21:03:57.7191923 +0800 CST m=+0.001547501
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:03:57.7186731 +0800 CST m=+0.001028301
	// layer:  1
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:03:57.7191923 +0800 CST m=+0.001547501
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: cccc, value: 111111, update_time: 2022-05-01 21:04:15.6276385 +0800 CST m=+0.001613801
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:04:15.627114 +0800 CST m=+0.001089301
	// layer:  2
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:04:15.6276385 +0800 CST m=+0.001613801
	cache.Add("aaaa", "111111")
	cache.OrderPrint()
	// 	layer:  0
	// key: cccc, value: 111111, update_time: 2022-05-01 21:04:33.6610861 +0800 CST m=+0.001573401
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:04:33.6605573 +0800 CST m=+0.001044601
	// layer:  3
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:04:33.6610861 +0800 CST m=+0.001573401
	cache.Reduce()
	// cache.Reduce()
	// cache.OrderPrint()
	// 	layer:  0
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:06:01.8192073 +0800 CST m=+0.002087801
	// key: cccc, value: 111111, update_time: 2022-05-01 21:06:01.8186799 +0800 CST m=+0.001560401
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:06:01.8181554 +0800 CST m=+0.00103590
	fmt.Println("-------------5555-------------")
	removekey, ok := cache.Add("bbb", "111111")
	fmt.Println("---------", ok)
	fmt.Println("---------", removekey)
	cache.OrderPrint()
	// 	layer:  0
	// key: bbb, value: 111111, update_time: 2022-05-01 21:25:33.6920657 +0800 CST m=+0.002121001
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:25:33.6915073 +0800 CST m=+0.001562601
	// key: cccc, value: 111111, update_time: 2022-05-01 21:25:33.6915073 +0800 CST m=+0.001562601
	cache.Add("eee", "111111")
	// 	cache.OrderPrint()
	// 	layer:  0
	// key: eee, value: 111111, update_time: 2022-05-01 21:39:09.9141285 +0800 CST m=+0.006492601
	// key: bbb, value: 111111, update_time: 2022-05-01 21:39:09.9097346 +0800 CST m=+0.002098701
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:39:09.9092087 +0800 CST m=+0.001572801

	if value, ok := cache.Get("cccc"); ok {
		fmt.Println(value) // "111111"
	}

}
