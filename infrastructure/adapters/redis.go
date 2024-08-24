package adapters

import (
	"github.com/redis/go-redis/v9"
	"github.com/wawayes/lark-bot/infrastructure"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(conf infrastructure.Config) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     conf.Redis.Addr,
			Password: conf.Redis.Password,
			DB:       conf.Redis.DB,
		}),
	}
}
