package reflect

import (
	"github.com/stretchr/testify/assert"
	"go-train/orm/reflect/types"
	"reflect"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			name:   "struct",
			entity: types.NewUser("justfong", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name:        "GetAge",
					InputTypes:  []reflect.Type{reflect.TypeOf(types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				//"ChangeName": {
				//	Name:       "ChangeName",
				//	InputTypes: []reflect.Type{reflect.TypeOf("")},
				//OutputTypes: []reflect.Type{reflect.TypeOf(0)},
				//Result:      []any{18},
				//},
			},
		},
		{
			name:   "pointer",
			entity: types.NewUserPtr("justfong", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name:        "GetAge",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				"ChangeName": {
					Name:        "ChangeName",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					OutputTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IterateFunc(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRes, res)
		})
	}
}
