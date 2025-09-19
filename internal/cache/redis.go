package cache

import (
	"context"
	"time"

	wbfredis "github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type RedisCache struct {
	client   *wbfredis.Client
	prefix   string
	retryStr retry.Strategy
}

func NewRedisCache(client *wbfredis.Client, prefix string, rs retry.Strategy) *RedisCache {
	return &RedisCache{
		client:   client,
		prefix:   prefix,
		retryStr: rs,
	}
}

func (r *RedisCache) key(k string) string {
	return r.prefix + k
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if r.retryStr.Attempts > 0 {
		return r.client.GetWithRetry(ctx, r.retryStr, r.key(key))
	}
	return r.client.Get(ctx, r.key(key))
}

func (r *RedisCache) Set(ctx context.Context, key, val string, ttl time.Duration) error {
	if r.retryStr.Attempts > 0 {
		delay := r.retryStr.Delay
		var lastErr error
		for i := 0; i < r.retryStr.Attempts; i++ {
			if err := r.client.Client.Set(ctx, r.key(key), val, ttl).Err(); err != nil {
				return nil
			} else {
				lastErr = err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * r.retryStr.Backoff)
		}
		return lastErr
	}
	return r.client.Client.Set(ctx, r.key(key), val, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	_, err := r.client.Client.Del(ctx, r.key(key)).Result()
	return err
}
