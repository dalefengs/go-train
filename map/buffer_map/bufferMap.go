package main

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	// [] 保存两个 map
	// []map[string]Data 排行榜单
	MapCache       = [2][]map[string]RankData{}
	MapIndex int32 = -1
)

type RankData struct {
	OpenId string `json:"openId"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
}

func callPlayerRank() {
	resultData := data
	nowIndex := atomic.LoadInt32(&MapIndex)
	if nowIndex < 0 {
		MapCache[0], MapCache[1] = resultData, resultData
		atomic.StoreInt32(&MapIndex, 0)
	} else if nowIndex == 0 {
		MapCache[1] = resultData
		atomic.StoreInt32(&MapIndex, 1)
	} else {
		MapCache[0] = resultData
		atomic.StoreInt32(&MapIndex, 0)
	}
	//fmt.Printf("map 长度 readMapLen = %d, writeMapLen = %d MapIndex = %d \n", len(MapCache[0]), len(MapCache[0]), MapIndex)
}

// 模拟一下数据
func mockData(count int) (ret []map[string]RankData) {
	for i := 0; i < count; i++ {
		ret = append(ret, map[string]RankData{
			strconv.Itoa(rand.Int()): RankData{
				Name:   "name",
				Score:  i,
				OpenId: strconv.Itoa(rand.Int()),
			},
		})
	}
	return
}

var data []map[string]RankData

func main() {
	data = mockData(2000)
	// 启动永久协程每1秒写入一次数据
	go func() {
		for {
			callPlayerRank()
			//time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
		}
	}()

	// 这里测试一直读取数据
	go func() {
		for {
			if MapIndex == -1 {
				continue
			}
			ret := MapCache[MapIndex]
			//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
			//time.Sleep(1 * time.Second)
			//time.Sleep(50 * time.Microsecond)
			go func() {
				for k, v := range ret {
					_ = v
					go func(data map[string]RankData) {
						for kk, vv := range data {
							_ = kk
							_ = vv
						}
					}(ret[k])
				}
			}()
		}
	}()
	go func() {
		for {
			if MapIndex == -1 {
				continue
			}
			ret := MapCache[MapIndex]
			//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
			//time.Sleep(1 * time.Second)
			//time.Sleep(50 * time.Microsecond)
			go func() {
				for k, v := range ret {
					_ = v
					go func(data map[string]RankData) {
						for kk, vv := range data {
							_ = kk
							_ = vv
						}
					}(ret[k])
				}
			}()
		}
	}()
	for {
		if MapIndex == -1 {
			continue
		}
		ret := MapCache[MapIndex]
		//fmt.Printf("key%d, Read value: %v \n", count, val["key"+strconv.Itoa(count)])
		//time.Sleep(1 * time.Second)
		//time.Sleep(50 * time.Microsecond)
		go func() {
			for k, v := range ret {
				_ = v
				go func(data map[string]RankData) {
					for kk, vv := range data {
						_ = kk
						_ = vv
					}
				}(ret[k])
			}
		}()
	}
}
