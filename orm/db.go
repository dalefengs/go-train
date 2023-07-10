package orm

import "reflect"

type DBOption func(db *DB)

type DB struct {
	r *registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	res := &DB{
		r: &registry{
			models: make(map[reflect.Type]*model, 64),
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
