package valuer

import "database/sql"

type Valuer interface {
	SetColumns(rows *sql.Rows) error
}

type Creator func(entity any) Valuer
