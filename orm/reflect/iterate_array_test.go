package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		want    []any
		wantErr error
	}{
		{
			name:   "slice",
			entity: []any{"li", "feng"},

			want:    []any{"li", "feng"},
			wantErr: nil,
		},
		{
			name:   "array",
			entity: [2]any{"li", "feng"},

			want:    []any{"li", "feng"},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateArray(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, res)
		})
	}
}

func TestIterateMap(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantKeys   []any
		wantValues []any
		wantErr    error
	}{
		{
			name: "slice",
			entity: map[any]any{
				"li":   18,
				"feng": 19,
			},

			wantKeys:   []any{"li", "feng"},
			wantValues: []any{18, 19},
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resKeys, resValues, err := IterateMap(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantKeys, resKeys)
			assert.Equal(t, tc.wantValues, resValues)
		})
	}
}
