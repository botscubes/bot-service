package redisauth

// two redis is bad practice

import (
	"github.com/botscubes/bot-service/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(c *config.RedisAuthConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Host + ":" + c.Port,
		Password: c.Pass,
		DB:       c.Db,
	})
}
