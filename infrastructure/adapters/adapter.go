package adapters

import (
	"github.com/wawayes/lark-bot/infrastructure"
)

type Adapter interface {
	Logger() *Logger
	Redis() *RedisClient
	Lark() *LarkClient
}

type AdapterImpl struct {
	logger *Logger
	redis  *RedisClient
	lark   *LarkClient
}

func NewAdapter(conf infrastructure.Config) *AdapterImpl {
	return &AdapterImpl{
		logger: NewLogger(conf),
		redis:  NewRedisClient(conf),
		lark:   NewLarkClient(conf),
	}
}

func (a *AdapterImpl) Logger() *Logger {
	return a.logger
}

func (a *AdapterImpl) Redis() *RedisClient {
	return a.redis
}

func (a *AdapterImpl) Lark() *LarkClient {
	return a.lark
}
