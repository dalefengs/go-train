package reflect

import "reflect"

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)
	//for typ.Kind() == reflect.Pointer {
	//	typ = typ.Elem()
	//}
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func
		numIn := fn.Type().NumIn()
		numOut := fn.Type().NumOut()

		input := make([]reflect.Type, 0, numIn)
		inputValue := make([]reflect.Value, 0, numIn)
		output := make([]reflect.Type, 0, numOut)

		inputValue = append(inputValue, reflect.ValueOf(entity))
		input = append(input, reflect.TypeOf(entity))

		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			input = append(input, fnInType)
			inputValue = append(inputValue, reflect.Zero(fnInType))
		}
		for j := 0; j < numOut; j++ {
			output = append(output, fn.Type().Out(j))
		}
		resValues := fn.Call(inputValue)
		result := make([]any, 0, len(resValues))
		for _, v := range resValues {
			result = append(result, v.Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
	}
	return res, nil

}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
