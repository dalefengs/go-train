package orm

import (
	"context"
	"go-train/orm/internal/errs"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	model *model
	sb    *strings.Builder
	args  []any
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
	s.sb = &strings.Builder{}

	m, err := parseModel(new(T))
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
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
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
		fd, ok := s.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnkonownField(exp.name)
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
