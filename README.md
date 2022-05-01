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
 > example
  ```go
package main

import (
	"fmt"

	"github.com/hyahm/gocache"
)

func main() {
	cache := gocache.NewCache[string, string](3, gocache.LRU)

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
```
