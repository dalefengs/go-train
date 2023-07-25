package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// 装饰器模式

var (
	errOverCapacity = errors.New("cache: 超过最大容量")
)

type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(c *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: c,
		maxCnt:          maxCnt,
	}
	origin := c.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	//cnt := atomic.AddInt32(&c.cnt, 1)
	//if cnt > c.maxCnt {
	//	atomic.AddInt32(&c.cnt, -1)
	//	return errOverCapacity
	//}
	//return c.BuildInMapCache.Set(ctx, key, val, expiration)
	c.mutex.Lock()
	_, ok := c.data[key]
	if !ok {
		c.cnt++
		if c.cnt > c.maxCnt {
			c.mutex.Unlock()
			return errOverCapacity
		}
	}
	c.mutex.Unlock()
	return c.BuildInMapCache.Set(ctx, key, val, expiration)
}
