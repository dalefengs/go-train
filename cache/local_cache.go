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

// BuildInMapCache 高性能的本地缓存
type BuildInMapCache struct {
	data      map[string]*item
	mutex     sync.RWMutex
	close     chan struct{}
	onEvicted func(key string, val any)
}

type BuildInMapCacheOption func(cache *BuildInMapCache)

func NewBuildInMapCache(interval time.Duration, opts ...BuildInMapCacheOption) *BuildInMapCache {
	res := &BuildInMapCache{
		data:      make(map[string]*item, 100),
		close:     make(chan struct{}),
		onEvicted: func(key string, val any) {},
	}
	for _, opt := range opts {
		opt(res)
	}

	// 启动清理过期的 goroutine
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
							res.delete(key)
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

// WithEvictedCallback 设置一个回调函数，当缓存条目被驱逐时调用该函数
func WithEvictedCallback(fn func(key string, val any)) BuildInMapCacheOption {
	return func(cache *BuildInMapCache) {
		cache.onEvicted = fn
	}
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

func (b *BuildInMapCache) delete(key string) {
	itm, ok := b.data[key]
	if !ok {
		return
	}
	delete(b.data, key)
	b.onEvicted(key, itm)
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
	return !i.deadline.IsZero() && i.deadline.Before(t)
}
