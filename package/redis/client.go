package redis_engine

import (
	"github.com/phamduytien1805/package/config"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

var RedisConn *redis.Client

func NewRedis(config *config.Config) *redis.Client {
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password, // no password set
		DB:       config.Redis.DB,       // use default DB
	})
	return RedisConn
}

func NewRedisStore(client *redis.Client) RedisQuerier {
	return &RedisStore{client}
}
