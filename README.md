# Build-in Go Cache（lru, lfu, alfu）
Thread-safe go language general simple, lru, lfu, alfu) algorithm

### What is alfu
The alfu algorithm is based on the lfu algorithm by adding a dynamic reduction level to complement the shortcomings of lfu.
add goroutine below
```go
	tick := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-tick.C:
			lfu.mu.Lock()
			for index, lru := range lfu.layer {
				if index == lfu.min {
					continue
				}
				key, value, update_time := lru.GetLastKeyUpdateTime()
				if time.Since(update_time).Hours() >= 24 {
					lru.Remove(key)
					lfu.cache[key] = index / 2
					newLevel := lfu.cache[key] / lfu.claddingSize
					if lfu.min > newLevel {
						lfu.min = newLevel
					}

					lfu.add(newLevel, key, value)

				}
			}
			lfu.mu.Unlock()
		}

	}
```
### Install (go version >= 1.18)
```
go get github.com/hyahm/gocache
```

### Cacher 
Usually only need to use `Add`, `Get` to meet most needs  
Add(key K, value V) (K, bool) // add key and value， If the number exceeds the maximum number, automatically delete the first eliminated value, and return the eliminated key and true, otherwise return the currently inserted key and false    
Remove(key K)                 // remove key      
Len() int                     // cache length    
OrderPrint()                  // order print key and value, if use simple,lru, last one eliminated first , if lfu, alfu, The first layer last is eliminated first   
Get(key K) (V, bool)          // get value by key  
LastKey() K                   // get key of last one  
Resize(int)                   // resize Max size  
### 
Before use, the number of caches must be initialized. If the memory is large enough, it can be set to a large size, and the efficiency will not be affected by the length.
If it exceeds the set value, it will automatically delete the value at the end, if it exists, it will automatically update this value to the beginning, and update the value
 > initialization
  ```
  // If the number of accesses to each layer follows the step size of 10 ALU and the number of accesses to each layer follows the algorithm of 1 ALU, then the number of accesses to each layer follows the step size of 10 ALU
  gocache.NewCache[comparable, any](<Max caches of number>, <Algorithm >, [Step size per layer])
  ```

> example lru
```go
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

```
> example lfu
```go
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
```
> example alfu
```go
package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	// only ALFU have Reduce function
	cache := gocache.NewCache[string, string](3, gocache.ALFU).(*gocache.Alfu[string, string])
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
	cache.Reduce()
	// cache.OrderPrint()
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
	// cache.OrderPrint()
	// 	layer:  0
	// key: cccc, value: 111111, update_time: 2022-05-01 21:04:33.6610861 +0800 CST m=+0.001573401
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:04:33.6605573 +0800 CST m=+0.001044601
	// layer:  3
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:04:33.6610861 +0800 CST m=+0.001573401
	cache.Reduce()
	cache.Reduce()
	// cache.OrderPrint()
	// 	layer:  0
	// key: aaaa, value: 111111, update_time: 2022-05-01 21:06:01.8192073 +0800 CST m=+0.002087801
	// key: cccc, value: 111111, update_time: 2022-05-01 21:06:01.8186799 +0800 CST m=+0.001560401
	// key: dddd, value: bbbbb, update_time: 2022-05-01 21:06:01.8181554 +0800 CST m=+0.00103590
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

```