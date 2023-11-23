package simplecache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	c *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		c: client,
	}
}

func (c *RedisClient) Get(
	ctx context.Context,
	key string,
) ([]byte, bool, error) {
	bytes, err := c.c.Get(ctx, key).Bytes()
	// キャッシュが存在しない場合
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get from redis: %w", err)
	}
	return bytes, true, nil
}

func (c *RedisClient) Keys(
	ctx context.Context,
	pattern string,
	bytes []byte,
	expiration time.Duration,
) ([]string, error) {
	cmd := c.c.Keys(ctx, pattern)

	return cmd.Result()
}

func (c *RedisClient) Set(
	ctx context.Context,
	key string,
	bytes []byte,
	expiration time.Duration,
) error {
	err := c.c.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set to redis: %w", err)
	}
	return nil
}

func (c *RedisClient) Del(
	ctx context.Context,
	key string,
) error {
	err := c.c.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to set to redis: %w", err)
	}
	return nil
}
