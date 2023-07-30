package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
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

type ConcurrentMap struct {
	readMap  unsafe.Pointer // 只读
	writeMap unsafe.Pointer // 只写
	once     sync.Once
}

func NewConcurrentMap(cap int) {
	if BufferMap == nil {
		readMap := make(map[string]Data, cap)
		writeMap := make(map[string]Data, cap)
		m := ConcurrentMap{
			readMap:  unsafe.Pointer(&readMap),
			writeMap: unsafe.Pointer(&writeMap),
		}
		BufferMap = &m
	}
	fmt.Printf("readMap %p, writeMap %p \n", *(*map[string]Data)(BufferMap.readMap), *(*map[string]Data)(BufferMap.writeMap))
}

func (m *ConcurrentMap) GetReadMap() map[string]Data {
	ret := *(*map[string]Data)(m.readMap)
	return ret
}

func (m *ConcurrentMap) GetWriteMap() map[string]Data {
	ret := *(*map[string]Data)(m.writeMap)
	return ret
}

func (m *ConcurrentMap) Set(key string, val Data) {
	newWriterMap := m.deepCopyMap(m.GetWriteMap())
	newWriterMap[key] = val
	m.CAS(&m.writeMap, m.writeMap, unsafe.Pointer(&newWriterMap))
	m.SwitchMapPointer()

}

func (m *ConcurrentMap) CAS(addr *unsafe.Pointer, oldPtr unsafe.Pointer, newPtr unsafe.Pointer) {
	for {
		swapped := atomic.CompareAndSwapPointer(addr, oldPtr, newPtr)
		if swapped {
			break
		}
	}
}

func (m *ConcurrentMap) Append(vals []Data) {
	writeMap := m.GetWriteMap()
	newWriterMap := m.deepCopyMap(writeMap)
	for _, v := range vals {
		newWriterMap[v.OpenId] = v
	}
	// 自旋
	m.CAS(&m.writeMap, m.writeMap, unsafe.Pointer(&newWriterMap))
}

// SwitchMapPointer 切换指针
func (m *ConcurrentMap) SwitchMapPointer() {
	oldWriteMap := m.writeMap
	// 将 readMap 指向 writeMap
	// writeMap 修改时会重新修改为深拷贝，所以不会影响 readMap
	m.CAS(&m.readMap, m.readMap, oldWriteMap)

}

func (m *ConcurrentMap) deepCopyMap(data map[string]Data) map[string]Data {
	newMap := make(map[string]Data, len(data))
	for k, v := range data {
		newMap[k] = v
	}
	return newMap
}

func mockData(count int) (ret []Data) {

	for i := 0; i < count; i++ {
		ret = append(ret, Data{
			Name:   "name",
			Score:  i,
			OpenId: strconv.Itoa(rand.Int()),
		})
	}
	time.Sleep(5 * time.Second)
	return
}

func main() {
	NewConcurrentMap(100000)
	// 启动永久协程每10秒写入一次数据
	go func() {
		for {
			for i := 0; i < 5; i++ {
				ret := mockData(500)
				s := time.Now()
				BufferMap.Append(ret)
				e := time.Now()
				cost := e.Sub(s).Milliseconds()
				fmt.Printf("Append %d ms \n", cost)
			}
			BufferMap.SwitchMapPointer()
		}
	}()

	//count := 1
	//go func() {
	//	for {
	//		//time.Sleep(1 * time.Second)
	//		d := Data{
	//			Name:  "test",
	//			Score: count,
	//		}
	//		BufferMap.Set("key"+strconv.Itoa(count), d)
	//		if count == 100000 {
	//			count = 1
	//		}
	//		count++
	//	}
	//}()

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
