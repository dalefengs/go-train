package reflect

import "reflect"

func IterateArray(entity any) ([]any, error) {
	res := make([]any, 0)
	val := reflect.ValueOf(entity)
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil
}

func IterateMap(entity any) ([]any, []any, error) {
	resKeys := make([]any, 0)
	resvals := make([]any, 0)

	val := reflect.ValueOf(entity)
	keys := val.MapKeys()
	for _, key := range keys {
		v := val.MapIndex(key)
		resKeys = append(resKeys, key.Interface())
		resvals = append(resvals, v.Interface())
	}

	return resKeys, resvals, nil
}
