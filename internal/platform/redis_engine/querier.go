package redis_engine

import (
	"context"
	"time"
)

type RedisQuerier interface {
	Set(ctx context.Context, key string, value interface{}) error
	SetTx(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetRaw(ctx context.Context, key string) (string, error)
	Get(ctx context.Context, key string, val interface{}) error
	Exist(ctx context.Context, key string) (bool, error)
}

func (c *RedisStore) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, key, value, 0).Err()
}

func (c *RedisStore) SetTx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}
func (c *RedisStore) GetRaw(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisStore) Get(ctx context.Context, key string, val interface{}) error {
	return c.client.Get(ctx, key).Scan(val)
}

func (c *RedisStore) Exist(ctx context.Context, key string) (bool, error) {
	val, err := c.client.Exists(ctx, key).Result()
	return val != 0, err
}
