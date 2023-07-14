package orm

type op string

const (
	opLT  = "<"
	opGT  = ">"
	opEq  = "="
	opOr  = "OR"
	opNot = "NOT"
	opAnd = "AND"
	opIn  = "IN"
)

func (o op) String() string {
	return string(o)
}

// Predicate 查询条件
type Predicate struct {
	left  Expression
	op    op
	right Expression
}

// Column 列名
type Column struct {
	name string
}

// value  参数列表
type value struct {
	val any
}

func (left Column) Eq(arg any) Predicate {
	return Predicate{
		left:  left,
		op:    opEq,
		right: value{val: arg},
	}
}

// LT <
func (left Column) LT(arg any) Predicate {
	return Predicate{
		left:  left,
		op:    opLT,
		right: value{val: arg},
	}
}

func Not(right Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: right,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

// Expression 是一个标记接口, 代表表达式
type Expression interface {
	expr()
}

func (Column) expr() {}

func (Predicate) expr() {}

func (value) expr() {}
