package unsafa

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnSafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnSafeAccessor(entity any) *UnSafeAccessor {
	typ := reflect.TypeOf(entity)
	numField := typ.NumField()
	field := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		field[fd.Name] = FieldMeta{
			Offset: fd.Offset,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnSafeAccessor{
		fields:  map[string]FieldMeta{},
		address: val.UnsafePointer(),
	}
}

func (a *UnSafeAccessor) Field(field string) (any, error) {
	// 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := uintptr(a.address) + fd.Offset

}

type FieldMeta struct {
	Offset uintptr
}
