package redis_engine

import (
	"github.com/phamduytien1805/package/config"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

var redisConn *redis.Client

func NewRedis(config *config.Config) *redis.Client {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password, // no password set
		DB:       config.Redis.DB,       // use default DB
	})
	return redisConn
}

func NewRedisStore(client *redis.Client) RedisQuerier {
	return &RedisStore{client}
}
