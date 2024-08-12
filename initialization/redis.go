package initialization

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func InitRedis(conf Config) {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     conf.Redis.Addr,
			Password: conf.Redis.Password,
			DB:       conf.Redis.DB,
		})
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}
