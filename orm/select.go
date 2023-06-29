package orm

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Build() (*Query, error) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	// 没有调用 FROM 就 使用反射获取表名称
	if s.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		sb.WriteByte('`')
		sb.WriteString(typ.Name())
		sb.WriteByte('`')

	} else {
		sb.WriteString(s.table)
	}
	sb.WriteByte(';')
	return &Query{
		Sql: sb.String(),
	}, nil

}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}