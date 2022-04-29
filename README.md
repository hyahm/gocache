# 缓存算法（lru, lfu, alfu）
Thread-safe go language general simple, lru, lfu, alfu) algorithm

### What is alfu
The alfu algorithm is based on the lfu algorithm by adding a dynamic reduction level to complement the shortcomings of lfu.
add goroutine below
```
	tick := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-tick.C:
			lfu.mu.Lock()
			for index, lru := range lfu.row {
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
go get github.com/hyahm/lru
```
### 
Before use, the number of caches must be initialized. If the memory is large enough, it can be set to a large size, and the efficiency will not be affected by the length.
If it exceeds the set value, it will automatically delete the value at the end, if it exists, it will automatically update this value to the beginning, and update the value
 > initialization
  ```
  // If the number of accesses to each layer follows the step size of 10 ALU and the number of accesses to each layer follows the algorithm of 1 ALU, then the number of accesses to each layer follows the step size of 10 ALU
  gocache.NewCache[comparable, any](<Max caches of number>, <Algorithm >, [Step size per layer])
  ```
 > example
  ```
package main

import (
	"fmt"
	"github.com/hyahm/gocache"
)


func main() {
	cache = gocache.NewCache[string, any](3, gocache.LFU, 2)

	cache.Add("adsf", "bbbbb")
	cache.Add("cccc", "111111")
	golog.Info(lru.Len())
	cache.OrderPrint(0)
}
```