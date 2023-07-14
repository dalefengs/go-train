package orm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-train/orm/internal/errs"
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
				Sql:  "SELECT * FROM `test_model`;",
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
				Sql:  "SELECT * FROM `test_model`;",
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
				Sql:  "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		}, {
			name:    "where not",
			builder: (&Selector[TestModel]{}).Where(Not(Column{"Age"}.Eq(18))),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		}, {
			name:    "where and",
			builder: (&Selector[TestModel]{}).Where(Column{"FirstName"}.Eq("feng").And(Column{"Age"}.Eq(18))),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE (`first_name` = ?) AND (`age` = ?);",
				Args: []any{"feng", 18},
			},
		}, {
			name:    "where or",
			builder: (&Selector[TestModel]{}).Where(Column{"FirstName"}.Eq("feng").Or(Column{"Age"}.Eq(18))),
			wantQuery: &Query{
				Sql:  "SELECT * FROM `test_model` WHERE (`first_name` = ?) OR (`age` = ?);",
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

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	// query err
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	// no row err
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	// one data
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Feng", "18", "RenGui")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	// sacn err
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("asd", "Feng", "18", "RenGui")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name string
		s    *Selector[TestModel]

		wantErr error
		wantRes any
	}{
		{
			name:    "build sql invalid",
			s:       NewSelector[TestModel](db).Where(Column{"XXX"}.Eq(1)),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
		{
			name:    "query err",
			s:       NewSelector[TestModel](db).Where(Column{"Id"}.Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "query no row err",
			s:       NewSelector[TestModel](db).Where(Column{"Id"}.Eq(1)),
			wantErr: ErrNoRows,
		},
		{
			name:    "one data",
			s:       NewSelector[TestModel](db).Where(Column{"Id"}.Eq(1)),
			wantErr: nil,
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Feng",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "RenGui"},
			},
		},
		//{
		//	name:    "scan err",
		//	s:       NewSelector[TestModel](db).Where(Column{"Id"}.Eq(1)),
		//	wantErr: ErrNoRows,
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSelector_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Feng", "18", "RenGui")
	rows.AddRow("2", "Li", "19", "Lisi")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name string
		s    *Selector[TestModel]

		wantErr error
		wantRes any
	}{
		{
			name:    "multi data",
			s:       NewSelector[TestModel](db).Where(Column{"Id"}.Eq(1)),
			wantErr: nil,
			wantRes: []*TestModel{
				{
					Id:        1,
					FirstName: "Feng",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "RenGui"},
				},
				{
					Id:        2,
					FirstName: "Li",
					Age:       19,
					LastName:  &sql.NullString{Valid: true, String: "Lisi"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func memoryDB(t *testing.T) *DB {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	return db
}

type TestModel struct {
	Id        int
	FirstName string
	Age       int8 `orm:"column=age"`
	LastName  *sql.NullString
}
