package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"
)

var BufferMap *ConcurrentMap

type Data struct {
	OpenId string `json:"openId"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
}

// ConcurrentMap 并发安全的 map
type ConcurrentMap struct {
	readMapPtr  unsafe.Pointer // 只读 Map 指针
	writeMapPtr unsafe.Pointer // 只写 Map 指针
}

func InitConcurrentMap(mapCap int) {
	if BufferMap != nil {
		return
	}
	readMap := make(map[string]Data, mapCap)
	writeMap := make(map[string]Data, mapCap)
	m := ConcurrentMap{
		readMapPtr:  unsafe.Pointer(&readMap),
		writeMapPtr: unsafe.Pointer(&writeMap),
	}
	BufferMap = &m
	fmt.Printf("readMapPtr %p, writeMapPtr %p \n", *(*map[string]Data)(BufferMap.readMapPtr), *(*map[string]Data)(BufferMap.writeMapPtr))
}

// LoadReadMap 获取只读 map
func (m *ConcurrentMap) LoadReadMap() map[string]Data {
	return *(*map[string]Data)(m.readMapPtr)
}

// LoadWriteMap 获取只写 map
func (m *ConcurrentMap) LoadWriteMap() map[string]Data {
	return *(*map[string]Data)(m.writeMapPtr)
}

// Set 设置写 Map 数据
func (m *ConcurrentMap) Set(key string, val Data) {
	newWriterMap := m.deepCopyMap(m.LoadWriteMap())
	newWriterMap[key] = val
	// 此时
	m.CAS(&m.writeMapPtr, m.writeMapPtr, unsafe.Pointer(&newWriterMap))
	m.SwitchReadMapPointer()
}

// CAS 原子操作 比较并交换指针
func (m *ConcurrentMap) CAS(addr *unsafe.Pointer, oldPtr unsafe.Pointer, newPtr unsafe.Pointer) {
	// 自旋锁 不断尝试 CAS 直到成功为止
	for {
		swapped := atomic.CompareAndSwapPointer(addr, oldPtr, newPtr)
		if swapped {
			break
		}
	}
}

// Append 追加一组数据
func (m *ConcurrentMap) Append(values []Data) {
	writeMap := m.LoadWriteMap()
	newWriterMap := m.deepCopyMap(writeMap)
	for _, v := range values {
		newWriterMap[v.OpenId] = v
	}
	m.CAS(&m.writeMapPtr, m.writeMapPtr, unsafe.Pointer(&newWriterMap))
}

// SwitchReadMapPointer 切换读Map指针到写Map
func (m *ConcurrentMap) SwitchReadMapPointer() {
	// 将 readMapPtr 指向 writeMapPtr
	// writeMapPtr 修改时会重新修改为深拷贝，所以不会影响 readMapPtr
	m.CAS(&m.readMapPtr, m.readMapPtr, m.writeMapPtr)
}

// deepCopyMap 深拷贝 Map, 防止切换后并发读写 Map
func (m *ConcurrentMap) deepCopyMap(data map[string]Data) map[string]Data {
	newMap := make(map[string]Data, len(data))
	for k, v := range data {
		newMap[k] = v
	}
	return newMap
}

// 模拟一下数据
func mockData(count int) (ret []Data) {
	for i := 0; i < count; i++ {
		ret = append(ret, Data{
			Name:   "name",
			Score:  i,
			OpenId: strconv.Itoa(rand.Int()),
		})
	}
	return
}

func main() {
	InitConcurrentMap(100000)
	// 启动永久协程每1秒写入一次数据
	go func() {
		for {
			for i := 0; i < 5; i++ {
				ret := mockData(500)
				s := time.Now()
				BufferMap.Append(ret)
				e := time.Now()
				cost := e.Sub(s).Milliseconds()
				time.Sleep(1 * time.Second)
				fmt.Printf("Append %d ms \n", cost)
			}
			BufferMap.SwitchReadMapPointer()
		}
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			fmt.Printf("map 长度 readMapLen = %d, writeMapLen = %d\n", len(BufferMap.LoadReadMap()), len(BufferMap.LoadWriteMap()))
		}
	}()

	// 这里测试一直读取数据
	go func() {
		for {
			ret := BufferMap.LoadReadMap()
			//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
			//time.Sleep(1 * time.Second)
			//time.Sleep(50 * time.Microsecond)
			go func() {
				for k, v := range ret {
					_ = k
					_ = v
				}
			}()
		}
	}()
	go func() {
		for {
			ret := BufferMap.LoadReadMap()
			//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
			//time.Sleep(1 * time.Second)
			//time.Sleep(50 * time.Microsecond)
			go func() {
				for k, v := range ret {
					_ = k
					_ = v
				}
			}()
		}
	}()
	for {
		ret := BufferMap.LoadReadMap()
		//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
		//time.Sleep(1 * time.Second)
		//time.Sleep(50 * time.Microsecond)
		go func() {
			for k, v := range ret {
				_ = k
				_ = v
			}
		}()
	}
}
