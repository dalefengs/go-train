package unsafa

import "testing"

func TestPrintFieldOffset(t *testing.T) {
	testCases := []struct {
		name   string
		entity any
	}{
		{
			name:   "user",
			entity: User{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			PrintFieldOffset(tc.entity)
		})
	}
}

type User struct {
	Name    string
	Age     int32
	Alias   []string
	Address string
}
