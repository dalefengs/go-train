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
	typ = typ.Elem()
	numField := typ.NumField()
	fields := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{
			Offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnSafeAccessor{
		fields:  fields,
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
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)
	// 知道确切类型
	//return *(*int)(fdAddress), nil

	// 不知道确切类型
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil
}

func (a *UnSafeAccessor) SetField(field string, val any) error {
	// 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}

type FieldMeta struct {
	Offset uintptr
	typ    reflect.Type
}
