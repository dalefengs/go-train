package main

import (
	"fmt"
	"time"
)

type ConcurrentMap struct {
	readMap  map[string]int
	writeMap map[string]int
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		readMap:  make(map[string]int),
		writeMap: make(map[string]int),
	}
}

func (cm *ConcurrentMap) Get(key string) (int, bool) {
	val, ok := cm.readMap[key]
	return val, ok
}

func (cm *ConcurrentMap) Set(key string, value int) {
	cm.writeMap[key] = value
}

func (cm *ConcurrentMap) Swap() {
	cm.readMap, cm.writeMap = cm.writeMap, make(map[string]int)
}

func main() {
	cm := NewConcurrentMap()

	// 启动永久协程每10秒写入一次数据
	go func() {
		for {
			cm.Set("key", 123) // 假设每次写入的值是123
			time.Sleep(10 * time.Second)
		}
	}()

	// 这里演示每隔1秒读取一次数据
	for {
		val, ok := cm.Get("key")
		if ok {
			fmt.Println("Read value:", val)
		} else {
			fmt.Println("Key not found")
		}
		cm.Swap() // 切换读写map
		time.Sleep(1 * time.Second)
	}
}
