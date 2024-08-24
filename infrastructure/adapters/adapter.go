package adapters

import "github.com/wawayes/lark-bot/infrastructure"

type Adapter interface {
	Redis() *RedisClient
	Lark() *LarkClient
}

type AdapterImpl struct {
	redis *RedisClient
	lark  *LarkClient
}

func NewAdapter(conf infrastructure.Config) *AdapterImpl {
	return &AdapterImpl{
		redis: NewRedisClient(conf),
		lark:  NewLarkClient(conf),
	}
}

func (a *AdapterImpl) Redis() *RedisClient {
	return a.redis
}

func (a *AdapterImpl) Lark() *LarkClient {
	return a.lark
}
