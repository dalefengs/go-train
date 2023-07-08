package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm：只支持一级指针结构体")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm：不支持的表达式 %v", expr)

}
func NewErrUnkonownField(name string) error {
	return fmt.Errorf("orm：未知字段 %v", name)

}
