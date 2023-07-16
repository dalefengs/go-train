package orm

import (
	"go-train/orm/internal/errs"
	"strings"
)

type builder struct {
	core
	sb   *strings.Builder
	args []any
}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	return b.buildExpression(p)
}

// 递归构建 Expression
func (b *builder) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	// 递归构建 Expression
	case Predicate:
		// 如果左右类型是 Predicate 那么会返回 true
		_, ok := exp.left.(Predicate)
		if ok {
			b.sb.WriteString("(")
		}
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			b.sb.WriteString(")")
		}

		// 操作符
		b.sb.WriteString(" " + exp.op.String())

		// 右边是表达式时 加上括号
		_, ok = exp.right.(Predicate)
		if ok {
			b.sb.WriteString(" (")
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}
	case Column:
		b.sb.WriteByte('`')
		fd, ok := b.model.fieldsMap[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		b.sb.WriteString(fd.colName)
		b.sb.WriteByte('`')
	case value:
		b.sb.WriteString(" ?")
		b.AddArg(exp.val)
	// 剩下不考虑
	default:
		//return errs.NewErrUnsupportedExpression(expr)
		return nil

	}
	return nil

}

func (b *builder) AddArg(val any) {
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, val)
}
