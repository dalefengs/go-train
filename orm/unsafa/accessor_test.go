package unsafa

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnSafeAccessor_Field(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	u := &User{Name: "Tom", Age: 18}
	accessor := NewUnSafeAccessor(u)
	field, err := accessor.Field("Age")
	require.NoError(t, err)
	assert.Equal(t, 18, field)

	err = accessor.SetField("Age", 19)
	require.NoError(t, err)
}
