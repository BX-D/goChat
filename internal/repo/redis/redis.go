package redis

import (
	"github.com/boxuanduan/gochat/config"
	"github.com/redis/go-redis/v9"
)

func Init(cfg *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: cfg.Addr, Password: cfg.Password})
}
