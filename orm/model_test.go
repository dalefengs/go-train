package orm

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-train/orm/internal/errs"
	"testing"
)

func Test_Register(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantModel any
		wantErr   error
		opts      []ModelOption
	}{
		{
			name:      "struct",
			entity:    TestModel{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fieldsMap: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
	}

	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Register(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, res)
		})
	}
}

func TestModelWithTableName(t *testing.T) {
	r := newRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName("test_model_ttt"))
	require.NoError(t, err) // 断定没有错误
	assert.Equal(t, "test_model_ttt", m.tableName)
}

func TestModelWithColumnName(t *testing.T) {
	testCases := []struct {
		name    string
		field   string
		colName string

		wantColName any
		wantErr     error
	}{
		{
			name:        "column name",
			field:       "FirstName",
			colName:     "AGE",
			wantColName: "AGE",
		},
		{
			name:    "invalid column name",
			field:   "xxxx",
			colName: "AGE",

			wantErr: errs.NewErrUnknownField("xxxx"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRegistry()
			register, err := r.Register(&TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := register.fieldsMap[tc.field]
			require.True(t, ok)
			assert.Equal(t, tc.wantColName, fd.colName)
		})
	}
}
