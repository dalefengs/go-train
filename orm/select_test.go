package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCases := []struct {
		name    string
		builder QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no from select",
			builder: &Selector[TestModel]{},
			wantQuery: &Query{
				Sql:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		}, {
			name:    "from select",
			builder: (&Selector[TestModel]{}).From("test_model"),
			wantQuery: &Query{
				Sql:  "SELECT * FROM test_model;",
				Args: nil,
			},
		}, {
			name:    "empty from select",
			builder: (&Selector[TestModel]{}).From(""),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		}, {
			name:    "db from select",
			builder: (&Selector[TestModel]{}).From("test_db.test_model"),
			wantQuery: &Query{
				Sql:  "SELECT * FROM test_db.test_model;",
				Args: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

type TestModel struct {
	ID        int
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
