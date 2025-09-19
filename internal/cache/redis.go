package cache

import (
	"context"
	"time"

	wbfredis "github.com/wb-go/wbf/redis"
	wbfretry "github.com/wb-go/wbf/retry"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type RedisCache struct {
	client   *wbfredis.Client
	prefix   string
	retryStr wbfretry.Strategy
}

func NewRedisCache(client *wbfredis.Client, prefix string, rs wbfretry.Strategy) *RedisCache {
	if rs.Attempts <= 0 {
		rs = internalRetry.DefaultStrategy
	}
	if rs.Delay <= 0 {
		rs.Delay = internalRetry.DefaultStrategy.Delay
	}
	if rs.Backoff <= 0 {
		rs.Backoff = internalRetry.DefaultStrategy.Backoff
	}
	return &RedisCache{client: client, prefix: prefix, retryStr: rs}
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
	attempts := r.retryStr.Attempts
	if attempts <= 0 {
		attempts = 1
	}

	delay := r.retryStr.Delay
	if delay <= 0 {
		delay = internalRetry.DefaultStrategy.Delay
	}

	backoff := r.retryStr.Backoff
	if backoff <= 0 {
		backoff = internalRetry.DefaultStrategy.Backoff
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := r.client.Client.Set(ctx, r.key(key), val, ttl).Err(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		time.Sleep(delay)
		delay = time.Duration(float64(delay) * backoff)
	}
	return lastErr
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	_, err := r.client.Client.Del(ctx, r.key(key)).Result()
	return err
}
