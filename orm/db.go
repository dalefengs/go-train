package orm

import (
	"sync"
)

type DBOption func(db *DB)

type DB struct {
	r *registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	res := &DB{
		r: &registry{
			models: sync.Map{},
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
