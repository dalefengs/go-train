package orm

import "strings"

type Deleter[T any] struct {
	builder
	table string
	where []Predicate
	model *model
	db    *DB
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sb = &strings.Builder{}
	m, err := d.db.r.parseModel(new(T))
	if err != nil {
		return nil, err
	}
	d.model = m
	sb := d.sb
	sb.WriteString("SELECT * FROM ")
	// 没有调用 FROM 就 使用反射获取表名称
	if d.table == "" {
		sb.WriteByte('`')
		sb.WriteString(d.model.tableName)
		sb.WriteByte('`')

	} else {
		sb.WriteString(d.table)
	}
	if len(d.where) > 0 {
		sb.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}
	}

	sb.WriteByte(';')
	return &Query{
		Sql:  sb.String(),
		Args: d.args,
	}, nil

}

func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

func (d *Deleter[T]) Where(predicate ...Predicate) *Deleter[T] {
	d.where = predicate
	return d
}
