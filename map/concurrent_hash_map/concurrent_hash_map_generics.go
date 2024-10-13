package concurrent_hash_map

import (
	"bytes"
	"encoding/binary"
	farmhash "github.com/leemcloughlin/gofarmhash"
	"sync"
	"unsafe"
)

// ConcurrentMap 支持高并发读写的Map
type ConcurrentMap[T comparable] struct {
	mps   []map[T]any    // 由多个小Map 组成的大Map
	seg   int            // 分段数
	locks []sync.RWMutex // 每一个小map独立的读写锁 。避免全局只有一把锁影响性能
	seed  uint32         // 每次执行 farmhash 传统一的 seed 种子，避免多线程冲突
}

func New[T comparable](seg int, cap int) *ConcurrentMap[T] {
	mps := make([]map[T]any, seg)
	locks := make([]sync.RWMutex, seg)
	for i := 0; i < seg; i++ {
		mps[i] = make(map[T]any, cap/seg)
		locks[i] = sync.RWMutex{}
	}
	return &ConcurrentMap[T]{
		mps:   mps,
		locks: locks,
		seg:   seg,
		seed:  0,
	}
}

// Pointer2Int 指针转 int
func Pointer2Int[T comparable](p *T) int {
	return *(*int)(unsafe.Pointer(p))
}

func IntToBytes(i int) []byte {
	buf := bytes.NewBuffer([]byte{})
	x := int64(i)
	// 使用大端序进行编码
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

// 判断 Key 写入到那个小map
func (m *ConcurrentMap[T]) getSegIndex(key T) int {
	hash := int(farmhash.Hash32WithSeed(IntToBytes(Pointer2Int(&key)), m.seed))
	return hash % m.seg
}

// Set 写入 key value
func (m *ConcurrentMap[T]) Set(key T, value any) {
	index := m.getSegIndex(key)
	m.locks[index].Lock()
	defer m.locks[index].Unlock()
	m.mps[index][key] = value
}

// Get 通过 key 获取 value
func (m *ConcurrentMap[T]) Get(key T) (any, bool) {
	index := m.getSegIndex(key)
	m.locks[index].RLock()
	defer m.locks[index].RUnlock()
	val, ok := m.mps[index][key]
	return val, ok
}
