package redis_engine

import (
	"github.com/phamduytien1805/package/config"
	"github.com/redis/go-redis/v9"
)

type RedisEngine struct {
	client *redis.Client
}

func NewRedis(config *config.Config) RedisQuerier {
	conn := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password, // no password set
		DB:       config.Redis.DB,       // use default DB
	})
	return &RedisEngine{
		client: conn,
	}

}
