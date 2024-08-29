package application

import (
	"context"
	"fmt"
	"time"

	"github.com/wawayes/lark-bot/infrastructure"
	"github.com/wawayes/lark-bot/infrastructure/adapters"
)

func InitBots(ctx context.Context, conf *infrastructure.Config, redisClient *adapters.RedisClient) error {
	for _, v := range conf.Bots {
		err := redisClient.Client.Set(ctx, fmt.Sprintf("bot:%s", v.BotName), v.ServiceFiled, time.Hour*24).Err()
		if err != nil {
			return fmt.Errorf("redis hset error: %w", err)
		}
	}
	return nil
}
