package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var BufferMap *ConcurrentMap

type Data struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ConcurrentMap struct {
	readMap atomic.Value // 只读
	once    sync.Once
}

func NewConcurrentMap(cap int) {
	if BufferMap == nil {
		m := ConcurrentMap{}
		m.readMap.Store(make(map[string]Data, cap))
		BufferMap = &m
	}
}

func (m *ConcurrentMap) GetReadMap() map[string]Data {
	return m.readMap.Load().(map[string]Data)
}

func (m *ConcurrentMap) Set(key string, val Data) {
	writeMap := m.readMap.Load().(map[string]Data)
	newWriteMap := make(map[string]Data, len(writeMap)-1)
	s := time.Now()
	// 将写 Map 中的数据复制到新的写 Map 中，但忽略要删除的键
	for k, v := range writeMap {
		newWriteMap[k] = v
	}
	e := time.Now()
	if e.Sub(s).Milliseconds() > 60 {
		fmt.Printf("for %v \n", e.Sub(s).Milliseconds())
	}
	newWriteMap[key] = val
	m.readMap.Store(newWriteMap)

	// 将写 Map 同步到读 Map
	m.readMap.Store(newWriteMap)
	m.SwitchMapPointer()
}

func (m *ConcurrentMap) SwitchMapPointer() {
}

func main() {
	NewConcurrentMap(10000000)
	// 启动永久协程每10秒写入一次数据
	count := 1
	go func() {
		for {
			//time.Sleep(1 * time.Second)
			d := Data{
				Name:  "test",
				Score: count,
			}
			BufferMap.Set("key"+strconv.Itoa(count), d)
			if count == 10000 {
				count = 1
			}
			count++
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Printf("readMapLen = %d\n", len(BufferMap.GetReadMap()))
		}
	}()

	// 这里演示每隔1秒读取一次数据
	for {
		s := time.Now()
		val := BufferMap.GetReadMap()
		e := time.Now()
		fmt.Printf("key%d, Read value: %v time: %v \n", count, val["key"+strconv.Itoa(count-1000)], e.Sub(s))
		time.Sleep(1 * time.Second)
	}
}
