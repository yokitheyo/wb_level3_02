package storage

import (
	"context"
	"time"

	wbfredis "github.com/wb-go/wbf/redis"
	wbfretry "github.com/wb-go/wbf/retry"
	internalRetry "github.com/yokitheyo/wb_level3_02/internal/retry"
)

type RedisCache struct {
	client   *wbfredis.Client
	prefix   string
	retryStr wbfretry.Strategy
}

func NewRedisCache(client *wbfredis.Client, prefix string, rs wbfretry.Strategy) Cache {
	if rs.Attempts <= 0 {
		rs = internalRetry.DefaultStrategy
	}
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
		return r.client.SetWithExpirationAndRetry(ctx, r.retryStr, r.key(key), val, ttl)
	}
	return r.client.SetWithExpiration(ctx, r.key(key), val, ttl)
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	if r.retryStr.Attempts > 0 {
		return r.client.DelWithRetry(ctx, r.retryStr, r.key(key))
	}
	return r.client.Del(ctx, r.key(key))
}
