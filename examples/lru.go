package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	cache := gocache.NewCache[string, string](3, gocache.LRU)
	cache.Add("dddd", "bbbbb")
	// cache.OrderPrint()
	// key: dddd, value: bbbbb, update_time: ...
	cache.Add("cccc", "111111")
	// cache.OrderPrint()
	// key: cccc, value: 111111, update_time: 2022-05-01 20:44:15.3180742 +0800 CST m=+0.001045101
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:44:15.3180742 +0800 CST m=+0.001045101
	cache.Add("cccc", "111111")
	// 	cache.OrderPrint()
	// 	key: cccc, value: 111111, update_time: 2022-05-01 20:44:27.4880738 +0800 CST m=+0.001043801
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:44:27.4880738 +0800 CST m=+0.001043801
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	key: aaaa, value: 111111, update_time: 2022-05-01 20:44:41.367326 +0800 CST m=+0.001033701
	// key: cccc, value: 111111, update_time: 2022-05-01 20:44:41.367326 +0800 CST m=+0.001033701
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:44:41.367326 +0800 CST m=+0.001033701
	cache.Add("aaaa", "111111")
	// 	cache.OrderPrint()
	// 	key: aaaa, value: 111111, update_time: 2022-05-01 20:44:41.367326 +0800 CST m=+0.001033701
	// key: cccc, value: 111111, update_time: 2022-05-01 20:44:41.367326 +0800 CST m=+0.001033701
	// key: dddd, value: bbbbb, update_time: 2022-05-01 20:44:41.367326 +0800 CST m=+0.001033701

	cache.Add("bbb", "111111")
	// 	cache.OrderPrint()
	// 	key: bbb, value: 111111, update_time: 2022-05-01 20:45:18.8995076 +0800 CST m=+0.001035101
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:45:18.8995076 +0800 CST m=+0.001035101
	// key: cccc, value: 111111, update_time: 2022-05-01 20:45:18.8995076 +0800 CST m=+0.001035101
	cache.Add("ddd", "111111")
	// 	cache.OrderPrint()
	// 	key: ddd, value: 111111, update_time: 2022-05-01 20:45:34.4605795 +0800 CST m=+0.000517201
	// key: bbb, value: 111111, update_time: 2022-05-01 20:45:34.4605795 +0800 CST m=+0.000517201
	// key: aaaa, value: 111111, update_time: 2022-05-01 20:45:34.4605795 +0800 CST m=+0.000517201

	if value, ok := cache.Get("cccc"); ok {
		fmt.Println(value) // "111111"
	}

}
