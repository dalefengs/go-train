package go_train

import (
	"fmt"
	"testing"
)

const YYYY_MM_DD = "2006-01-02"

type TestStruct struct {
	Name string
}

func TestFunc(t *testing.T) {
	arr := []string{"你好", "我很好", "3", "4", "11", "12", "13", "14", "111", "112", "113", "114"}
	for i, v := range arr {
		g := i / 4
		if g == 2 {
			fmt.Println(v)
		}
	}
}
