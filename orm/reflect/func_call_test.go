package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantRes map[string]FuncInfo
		wantErr error
	}{
		{},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IterateFunc(tt.entity)
			assert.Equal(t, tt.want, res)
		})
	}
}
