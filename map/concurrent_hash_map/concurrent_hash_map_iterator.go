package concurrent_hash_map

// 设计模式 迭代器模式

type MapIterator interface {
	Next() *MapEntity
}

type MapEntity struct {
	Key   string
	Value any
}

type ConcurrentHashMapIterator struct {
	m        *ConcurrentHashMap
	keys     [][]string
	rowIndex int
	colIndex int
}

func (iter *ConcurrentHashMapIterator) Next() *MapEntity {
	if iter.rowIndex >= len(iter.keys) {
		return nil
	}

	row := iter.keys[iter.rowIndex]
	if len(row) == 0 { // 本行为空 递归找下一行
		iter.rowIndex++
		return iter.Next()
	}
	key := row[iter.colIndex]

	value, _ := iter.m.Get(key)
	if iter.colIndex >= len(row)-1 {
		iter.rowIndex++
		iter.colIndex = 0
	} else {
		iter.colIndex++
	}

	return &MapEntity{
		Key:   key,
		Value: value,
	}
}
