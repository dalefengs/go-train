package concurrent_hash_map

import (
	farmhash "github.com/leemcloughlin/gofarmhash"
	"maps"
	"sync"
)

// ConcurrentHashMap 支持高并发读写的Map
type ConcurrentHashMap struct {
	mps   []map[string]any // 由多个小Map 组成的大Map
	seg   int              // 分段数
	locks []sync.RWMutex   // 每一个小map独立的读写锁 。避免全局只有一把锁影响性能
	seed  uint32           // 每次执行 farmhash 传统一的 seed 种子，避免多线程冲突
}

func NewConcurrentHashMap(seg int, cap int) *ConcurrentHashMap {
	mps := make([]map[string]any, seg)
	locks := make([]sync.RWMutex, seg)
	for i := 0; i < seg; i++ {
		mps[i] = make(map[string]any, cap/seg)
		locks[i] = sync.RWMutex{}
	}
	return &ConcurrentHashMap{
		mps:   mps,
		locks: locks,
		seg:   seg,
		seed:  0,
	}
}

// 判断 Key 写入到那个小map
func (m *ConcurrentHashMap) getSegIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed(IntToBytes(Pointer2Int(&key)), m.seed))
	return hash % m.seg
}

// Set 写入 key value
func (m *ConcurrentHashMap) Set(key string, value any) {
	index := m.getSegIndex(key)
	m.locks[index].Lock()
	defer m.locks[index].Unlock()
	m.mps[index][key] = value
}

// Get 通过 key 获取 value
func (m *ConcurrentHashMap) Get(key string) (any, bool) {
	index := m.getSegIndex(key)
	m.locks[index].RLock()
	defer m.locks[index].RUnlock()
	val, ok := m.mps[index][key]
	return val, ok
}

func (m *ConcurrentHashMap) CreateIterator() *ConcurrentHashMapIterator {
	keys := make([][]string, 0, len(m.mps))
	for _, mp := range m.mps {
		row := maps.Keys(mp)
		key := make([]string, 0, len(mp))
		row(func(s string) bool {
			key = append(key, s)
			return true
		})
		keys = append(keys, key)
	}
	return &ConcurrentHashMapIterator{
		m:        m,
		keys:     keys,
		rowIndex: 0,
		colIndex: 0,
	}
}
