package redis

import (
	"context"
	"errors"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/redis/go-redis/v9"
)

func (rdb *Rdb) SetComponents(botId int64, comps *[]*model.Component) error {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":component"

	for _, v := range *comps {
		if err := rdb.HSet(ctx, key, strconv.FormatInt(v.Id, 10), v).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (rdb *Rdb) GetComponent(botId int64, compId int64) (*model.Component, error) {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":component"

	component := &model.Component{}

	result, err := rdb.HGet(ctx, key, strconv.FormatInt(compId, 10)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if result == "" {
		return nil, errors.New("not found")
	}

	if err := component.UnmarshalBinary([]byte(result)); err != nil {
		return nil, err
	}

	return component, nil
}

func (rdb *Rdb) PrintAllComponents(botId int64) {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":component"
	log.Debug("key")
	log.Debug(key)

	data := rdb.HGetAll(ctx, key).Val()
	log.Debug(data)
}
