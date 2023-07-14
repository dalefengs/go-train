package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm：只支持一级指针结构体")
	ErrNoRows      = errors.New("orm：没有数据")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm：不支持的表达式 %v", expr)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm：未知字段 %v", name)
}

func NewErrUnknownColumn(name string) error {
	return fmt.Errorf("orm：未知列 %v", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm：非法标签 %v", pair)
}
