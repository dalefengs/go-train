package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var BufferMap *ConcurrentMap

type Data struct {
	Name   string `json:"name"`
	Score  int    `json:"score"`
	OpenId string `json:"openId"`
}

type ConcurrentMap struct {
	readMap  atomic.Value // 只读
	writeMap atomic.Value // 只写
	once     sync.Once
}

func NewConcurrentMap(cap int) {
	if BufferMap == nil {
		m := ConcurrentMap{}
		m.readMap.Store(make(map[string]Data, cap))
		m.writeMap.Store(make(map[string]Data, cap))
		BufferMap = &m
	}
}

func (m *ConcurrentMap) GetReadMap() map[string]Data {
	return m.readMap.Load().(map[string]Data)
}

func (m *ConcurrentMap) GetWriteMap() map[string]Data {
	return m.writeMap.Load().(map[string]Data)
}

func (m *ConcurrentMap) Append(vals []Data) {
	writeMap := m.writeMap.Load().(map[string]Data)
	newWriteMap := make(map[string]Data, len(writeMap)-1)
	for k, v := range writeMap {
		newWriteMap[k] = v
	}
	for _, v := range vals {
		newWriteMap[v.OpenId] = v
	}
	m.writeMap.Store(newWriteMap)
}

func (m *ConcurrentMap) AppendV1(vals []Data) {
	writeMap := m.writeMap.Load().(map[string]Data)
	for _, v := range vals {
		writeMap[v.OpenId] = v
	}
}

func (m *ConcurrentMap) Set(key string, val Data) {
	writeMap := m.GetWriteMap()
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
	m.writeMap.Store(newWriteMap)

	// 将写 Map 同步到读 Map
	m.readMap.Store(newWriteMap)
}

func (m *ConcurrentMap) SyncRWMap() {
	writeMap := m.GetWriteMap()
	m.readMap.Store(writeMap)
}

func mockData(count int) (ret []Data) {
	for i := 0; i < count; i++ {
		ret = append(ret, Data{
			Name:   "name",
			Score:  i,
			OpenId: strconv.Itoa(rand.Int()),
		})
	}
	time.Sleep(2 * time.Second)
	return
}

func main() {
	NewConcurrentMap(10000)
	// 启动永久协程每10秒写入一次数据
	go func() {
		for {
			for i := 0; i < 5; i++ {
				ret := mockData(500)
				s := time.Now()
				BufferMap.AppendV1(ret)
				e := time.Now()
				cost := e.Sub(s).Nanoseconds()
				fmt.Printf("AppendV1 %d ns \n", cost)
			}
			BufferMap.SyncRWMap()
		}
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			fmt.Printf("readMapLen = %d, writeMapLen = %d\n", len(BufferMap.GetReadMap()), len(BufferMap.GetWriteMap()))
		}
	}()

	// 这里演示每隔1秒读取一次数据
	go func() {
		for {

			ret := BufferMap.GetReadMap()
			//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
			//time.Sleep(1 * time.Second)
			//time.Sleep(50 * time.Microsecond)
			for k, v := range ret {
				_ = k
				_ = v
			}
		}
	}()
	go func() {
		for {

			ret := BufferMap.GetReadMap()
			//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
			//time.Sleep(1 * time.Second)
			//time.Sleep(50 * time.Microsecond)
			for k, v := range ret {
				_ = k
				_ = v
			}
		}
	}()
	for {

		ret := BufferMap.GetReadMap()
		//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
		//time.Sleep(1 * time.Second)
		//time.Sleep(50 * time.Microsecond)
		for k, v := range ret {
			_ = k
			_ = v
		}
	}
}
