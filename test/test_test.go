package main

import (
	"fmt"
	"sync"
	"testing"
)

type Data struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	m     map[string]Data
	once  *sync.Once
}

func TestOther(t *testing.T) {
	m := map[string]string{
		"1": "1",
	}
	b := m
	fmt.Printf("A%p  B%p\n", m, b)
}
