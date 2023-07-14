package redis

import (
	"errors"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("not found")
)

type Rdb struct {
	*redis.Client
}

func NewClient(c *config.RedisConfig) *Rdb {
	return &Rdb{
		redis.NewClient(&redis.Options{
			Addr:     c.Host + ":" + c.Port,
			Password: c.Pass,
			DB:       c.Db,
		}),
	}
}
