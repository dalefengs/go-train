package concurrent_hash_map

import (
	"math/rand"
	"sync"
	"testing"
)

var conMap = New[int64](8, 1000)

var syncMap = sync.Map{}

func readConMap() {
	for i := 0; i < 1000; i++ {
		key := rand.Int63()
		conMap.Get(key)
	}
}

func writerConMap() {
	for i := 0; i < 1000; i++ {
		key := rand.Int63()
		conMap.Set(key, 1)
	}
}

func readSyncMap() {
	for i := 0; i < 1000; i++ {
		key := rand.Int63()
		syncMap.Load(key)
	}
}

func writeSyncMap() {
	for i := 0; i < 1000; i++ {
		key := rand.Int63()
		syncMap.Store(key, 1)
	}
}

func BenchmarkConMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const P = 300
		wg := sync.WaitGroup{}
		wg.Add(2 * P)
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					readConMap()
				}
			}()
		}
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				writerConMap()
			}()
		}
	}
}

func BenchmarkSyncMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const P = 300
		wg := sync.WaitGroup{}
		wg.Add(2 * P)
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					readSyncMap()
				}
			}()
		}
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				writeSyncMap()
			}()
		}
	}
}
