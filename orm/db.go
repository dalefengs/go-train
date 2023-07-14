package orm

import (
	"database/sql"
	"sync"
)

type DBOption func(db *DB)

// DB 是 sql.DB 的装饰器
type DB struct {
	r  *registry
	db *sql.DB
}

func Open(driver, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r: &registry{
			models: sync.Map{},
		},
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil

}

func MustOpen(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		return nil
	}
	return db
}
