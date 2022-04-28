# 缓存算法（lru, lfu, alfu）
 线程安全的go语言通用simple, lru, lfu, alfu）算法,   
 alfu算法是在lfu算法基础增加动态减少层级来补全lfu的弊端   
### 安装 (go version >= 1.18)
```
go get github.com/hyahm/lru
```
### 使用

在使用前, 先要初始化缓存个数, 内存足够大的话, 可以设置很大, 不会因为长度影响效率, 
超过设定值会自动删除末尾的值, 如果存在的话会自动更新此值到开头, 更新值
 > 初始化(初始化完成后, 可以在任何地方调用方法)
  ```
  // 每层步长: lfu 或者 alfu 访问次数步长间距,  如果设置10，那么 1次访问和9次之间的访问都在这一层遵循 lru算法
  gocache.NewCache[comparable, any](<缓存个数>, <算法>, [每层步长])
  ```
 > 添加 key和value, 以下为doc中的例子
  ```
package main

import (
	"fmt"
	"lru"
)

type el struct {
	Id int
	Name string
}

func main() {
	cache = cache.New[string,string](10, cache.LRU)

	cache.Add("adsf", "bbbbb")
	cache.Add("cccc", "111111")
	golog.Info(lru.Len())
	cache.OrderPrint(0)
}
```
> 万能的add方法, 只要是添加值都可以使用此方法, 存在就会更新, 不存在就会插入
```
cache.Add(key, value any)
```
> 顺序打印(调试用)
```
cache.OrderPrint(level int)
```

> 删除key
```
cache.Remove(key any)
```
> 获取所有的key, 没有就返回空, 返回的key因为执行时间的问题, 可能导致有些key被删除了
```
cache.Keys(key any) []any
```
> 获取缓存长度 
```
cache.Len() uint64
```
> 根据key获取值
```
cache.Get(key any) any
```


基本上这些就能满足需求
