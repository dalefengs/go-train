package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	errKeyNotFound = errors.New("cache: 键不存在")
	errKeyExpired  = errors.New("cache: 键已过期")
)

type item struct {
	val      any
	deadline time.Time
}

type BuildInMapCache struct {
	data  map[string]*item
	mutex sync.RWMutex
	close chan struct{}
}

func NewBuildInMapCache(interval time.Duration) *BuildInMapCache {
	res := &BuildInMapCache{
		data: make(map[string]*item, 100),
	}
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for t := range ticker.C {
					res.mutex.Lock()
					i := 0
					for key, val := range res.data {
						if i > 1000 {
							break
						}
						// 设置了过期时间
						if !val.deadline.IsZero() && val.deadline.Before(t) {
							delete(res.data, key)
						}
						i++
					}
					res.mutex.Unlock()
				}
			case <-res.close:
				return
			}
		}
	}()
	return res
}

func (b *BuildInMapCache) Get(ctx context.Context, key string) (any, error) {
	b.mutex.RLock()
	res, ok := b.data[key]
	b.mutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w, key: %s", errKeyNotFound, key)
	}
	now := time.Now()
	// double chchek
	if res.deadlineBefore(now) {
		b.mutex.Lock()
		res, ok = b.data[key]
		if res.deadlineBefore(now) {
			delete(b.data, key)
		}
		b.mutex.Unlock()
		return nil, fmt.Errorf("%w key: %s", errKeyExpired, key)
	}
	return res.val, nil
}

func (b *BuildInMapCache) Set(ctx context.Context, key string, val any, expriation time.Duration) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var dl time.Time
	if expriation > 0 {
		dl = time.Now().Add(expriation)
	}
	b.data[key] = &item{
		val:      val,
		deadline: dl,
	}
	return nil
}

func (b *BuildInMapCache) Delete(ctx context.Context, key string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	delete(b.data, key)
	return nil
}

func (b *BuildInMapCache) Close(ctx context.Context, key string) error {
	select {
	case b.close <- struct{}{}:
	default:
		return errors.New("cache: 重复关闭")
	}
	return nil
}

// 是否过期
func (i *item) deadlineBefore(t time.Time) bool {
	return i.deadline.IsZero() && i.deadline.Before(t)
}
