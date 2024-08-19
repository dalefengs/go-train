package cache

import (
	"container/list"
	"fmt"
	"strconv"
	"testing"
)

// LRUCache基于最近最少使用（LRU）算法的缓存。
type LRUCache struct {
	cache map[int]string
	list  *list.List // 双端链表，用于辅助LRU决策
	cap   int
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cache: make(map[int]string, cap),
		list:  list.New(),
		cap:   cap,
	}
}

func (c *LRUCache) Add(key int, val string) {
	// 当cache容量不足时，淘汰末尾元素
	if c.list.Len() >= c.cap {
		back := c.list.Back()
		if back != nil {
			delete(c.cache, back.Value.(int))
			c.list.Remove(back)
		}
	}
	// 将新增元素写入链表头部
	c.list.PushFront(key)
	c.cache[key] = val
}

func (c *LRUCache) find(key int) *list.Element {
	head := c.list.Front()
	for head != nil {
		if head.Value.(int) == key {
			return head
		}
		head = head.Next()
	}
	return nil
}

func (c *LRUCache) Get(key int) (string, bool) {
	val, ok := c.cache[key]
	if !ok {
		return "", false
	}

	find := c.find(key)
	if find == nil {
		c.list.PushFront(key)
	} else {
		c.list.MoveToFront(find)
	}
	return val, true
}

func TestLRUCache(t *testing.T) {
	lru := NewLRUCache(10) // 缓存容量为10

	// 填满缓存
	for i := 0; i < 10; i++ {
		lru.Add(i, strconv.Itoa(i)) // 9876543210
	}

	// 访问偶数元素。被访问的元素会放到链表的首部
	for i := 0; i < 10; i += 2 {
		lru.Get(i) // 8642097531
	}

	// 再添加5个新元素。新添加的元素会放到链表的首部
	for i := 10; i < 15; i++ {
		lru.Add(i, strconv.Itoa(i)) // 14 13 12 11 10 8 6 4 2 0
	}

	// 检查缓存中还有没有最初的那10个元素
	for i := 0; i < 10; i++ {
		_, exists := lru.Get(i)
		fmt.Printf("key %d exists %t\n", i, exists) // 9 7 5 3 1不存在, 8 6 4 2 0 存在
	}
}
