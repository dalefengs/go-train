package orm

import (
	"context"
	"database/sql"
)

// 核心接口定义

// Querier 用于 Select 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor 用于 Insert, Update, Delete
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	Sql  string
	Args []any
}

type TableName interface {
	TableName() string
}
