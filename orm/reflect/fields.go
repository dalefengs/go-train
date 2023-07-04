package reflect

import (
	"errors"
	"reflect"
)

// IterateFields 遍历字段
func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}

	// 使用 for 的原因是可能是多重指针
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持类型")
	}

	fieldNums := typ.NumField()

	res := make(map[string]any, fieldNums)
	for i := 0; i < fieldNums; i++ {
		typName := typ.Field(i)
		if typName.IsExported() {
			res[typ.Field(i).Name] = val.Field(i).Interface()
		} else {
			res[typName.Name] = reflect.Zero(typName.Type).Interface()
		}
	}
	return res, nil
}

func SetField(entity any, field string, newValue any) (err error) {
	if entity == nil {
		return errors.New("不支持 nil")
	}
	val := reflect.ValueOf(entity)
	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("不支持的类型")
	}
	fieldVal := val.FieldByName(field)
	// 是否可以被修改
	if !fieldVal.CanSet() {
		return errors.New("不可修改字段")
	}
	fieldVal.Set(reflect.ValueOf(newValue))

	return nil
}
