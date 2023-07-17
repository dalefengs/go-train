package orm

import (
	"context"
	"database/sql"
	"go-train/orm/internal/errs"
	"reflect"
	"strings"
	"unsafe"
)

type Selector[T any] struct {
	builder
	table string
	where []Predicate
	db    *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		builder: builder{
			sb: &strings.Builder{},
		},
		db: db,
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	// 处理结果集
	db := s.db.db
	rows, err := db.QueryContext(ctx, q.Sql, q.Args...)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, ErrNoRows
	}
	// 获取到 select {} from 的字段
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	return s.parseRowsResult(rows, cs)
}

func (s *Selector[T]) parseRowsResult(rows *sql.Rows, columns []string) (*T, error) {
	vals := make([]any, 0, len(columns))
	valElems := make([]reflect.Value, 0, len(columns))
	for _, c := range columns {
		// c 是列名
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)
		}
		// 反射创建了一个新的实例
		val := reflect.New(fd.typ)
		vals = append(vals, val.Interface())
		valElems = append(valElems, val.Elem())
	}

	err := rows.Scan(vals...)
	if err != nil {
		return nil, err
	}
	tp := new(T)
	// 其实地址
	address := reflect.ValueOf(tp).UnsafePointer()
	for _, c := range columns {
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)
		}
		fdAddress := unsafe.Pointer(uintptr(address) + fd.offset)
		val := reflect.NewAt(fd.typ, fdAddress)
		vals = append(vals, val.Interface())
	}
	err = rows.Scan(vals)
	return tp, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	// 处理结果集
	db := s.db.db
	rows, err := db.QueryContext(ctx, q.Sql, q.Args...)
	result := make([]*T, 0)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		res, err := s.parseRowsResult(rows, columns)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	m, err := s.db.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	s.model = m
	sb := s.sb
	sb.WriteString("SELECT * FROM ")
	// 没有调用 FROM 就 使用反射获取表名称
	if s.table == "" {
		sb.WriteByte('`')
		sb.WriteString(s.model.tableName)
		sb.WriteByte('`')

	} else {
		sb.WriteString(s.table)
	}
	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		if err := s.buildPredicates(s.where); err != nil {
			return nil, err
		}
	}

	sb.WriteByte(';')
	return &Query{
		Sql:  sb.String(),
		Args: s.args,
	}, nil

}

// 递归构建 Expression
func (s *Selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	// 递归构建 Expression
	case Predicate:
		// 如果左右类型是 Predicate 那么会返回 true
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteString("(")
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteString(")")
		}

		// 操作符
		s.sb.WriteString(" " + exp.op.String())

		// 右边是表达式时 加上括号
		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteString(" (")
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:
		s.sb.WriteByte('`')
		fd, ok := s.model.fieldsMap[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteString(" ?")
		s.AddArg(exp.val)
	// 剩下不考虑
	default:
		//return errs.NewErrUnsupportedExpression(expr)
		return nil

	}
	return nil

}

func (s *Selector[T]) AddArg(val any) *Selector[T] {
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, val)
	return s
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}
