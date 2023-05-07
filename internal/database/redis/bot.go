package redis

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
)

func (rdb *Rdb) SetComponents(botId int64, comps *[]*model.Component) error {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":components"

	for _, v := range *comps {
		if err := rdb.HSet(ctx, key, strconv.FormatInt(v.Id, 10), v).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (rdb *Rdb) GetComponet(botId int64, compId int64) (*model.Component, error) {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":components"

	component := &model.Component{}

	result, err := rdb.HGet(ctx, key, strconv.FormatInt(compId, 10)).Result()
	log.Debug(err)

	if err := component.UnmarshalBinary([]byte(result)); err != nil {
		log.Debug(err)
	}

	return component, nil
}

func (rdb *Rdb) PrintAllComponents(botId int64) {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":components"

	data := rdb.HGetAll(ctx, key).Val()
	log.Debug(data)
}
