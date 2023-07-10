package orm

import (
	"go-train/orm/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string // 列名
}

type registry struct {
	lock   sync.RWMutex
	models map[reflect.Type]*model
}

func newRegistry() *registry {
	return &registry{
		models: make(map[reflect.Type]*model, 64),
	}
}

// double check lock
func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	r.lock.RLock()
	m, ok := r.models[typ]
	r.lock.RUnlock()
	if ok {
		return m, nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	m, ok = r.models[typ]
	if ok {
		return m, nil
	}
	var err error
	m, err = r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m
	return m, nil
}

// parseModel 解析模型
func (r *registry) parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	// 限制只能用一级指针
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numFields := typ.NumField()
	fieldMap := make(map[string]*field, numFields)
	for i := 0; i < numFields; i++ {
		fd := typ.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fields:    fieldMap,
	}, nil
}

// 大小写转换
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)

}
