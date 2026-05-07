// 浣滅敤锛歊edis 缂撳瓨閫傞厤灞傚疄鐜?
package redis

import (
	"context"
	"go-clean-architecture/config"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet 渚?Wire 浣跨敤鐨勪緷璧栨彁渚涜€呴泦鍚?
var ProviderSet = wire.NewSet(NewRedisClient, NewCacheClient)

// NewRedisClient 鏍规嵁閰嶇疆鍒濆鍖?Redis 杩炴帴
func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 妫€鏌ヨ繛閫氭€?
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

type CacheClient struct {
	client *redis.Client
}

// NewCacheClient 缂撳瓨瀹炵幇鏋勯€犲櫒
func NewCacheClient(client *redis.Client) *CacheClient {
	return &CacheClient{client: client}
}

func (c *CacheClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

