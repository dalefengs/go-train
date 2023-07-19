package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildInMapCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		cache   *BuildInMapCache
		wantVal any
		wantErr error
	}{
		{
			name: "key not found",
			key:  "not exist key",
			cache: func() *BuildInMapCache {
				return NewBuildInMapCache(10 * time.Second)
			}(),
			wantErr: fmt.Errorf("%w, key: %s", errKeyNotFound, "not exist key"),
		},
		{
			name: "get val",
			key:  "key1",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(10 * time.Second)
				err := res.Set(context.Background(), "key1", 123, time.Minute)
				require.NoError(t, err)
				return res
			}(),
			wantErr: nil,
			wantVal: 123,
		},
		{
			name: "expired",
			key:  "expired key",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(10 * time.Second)
				err := res.Set(context.Background(), "expired key", 123, 1*time.Second)
				require.NoError(t, err)
				time.Sleep(2 * time.Second)
				return res
			}(),
			wantErr: fmt.Errorf("%w key: %s", errKeyExpired, "expired key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}

}

func TestBuildInMapCache_Loop(t *testing.T) {
	cnt := 0
	c := NewBuildInMapCache(time.Second, WithEvictedCallback(func(key string, val any) {
		cnt++
	}))
	err := c.Set(context.Background(), "key1", 123, time.Second)
	require.NoError(t, err)
	time.Sleep(3 * time.Second)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data["key1"]
	require.False(t, ok)
	require.Equal(t, 1, cnt)

}
