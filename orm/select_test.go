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
		}, {
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(Column{"Age"}.Eq(18)),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `TestModel` WHERE `Age` = ?;",
				Args: []any{18},
			},
		}, {
			name:    "where not",
			builder: (&Selector[TestModel]{}).Where(Not(Column{"Age"}.Eq(18))),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `TestModel` WHERE  NOT (`Age` = ?);",
				Args: []any{18},
			},
		}, {
			name:    "where and",
			builder: (&Selector[TestModel]{}).Where(Column{"Name"}.Eq("feng").And(Column{"Age"}.Eq(18))),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `TestModel` WHERE (`Name` = ?) AND (`Age` = ?);",
				Args: []any{"feng", 18},
			},
		}, {
			name:    "where or",
			builder: (&Selector[TestModel]{}).Where(Column{"Name"}.Eq("feng").Or(Column{"Age"}.Eq(18))),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `TestModel` WHERE (`Name` = ?) OR (`Age` = ?);",
				Args: []any{"feng", 18},
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
