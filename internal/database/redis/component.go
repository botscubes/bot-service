package redis

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
)

func (rdb *Rdb) SetComponent(botId int64, comp *model.Component) error {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":component"

	if err := rdb.HSet(ctx, key, strconv.FormatInt(comp.Id, 10), comp).Err(); err != nil {
		return err
	}

	return rdb.Expire(ctx, key, config.RedisExpire).Err()
}

func (rdb *Rdb) DelComponent(botId int64, compId int64) error {
	ctx := context.Background()
	key := "bot" + strconv.FormatInt(botId, 10) + ":component"
	return rdb.HDel(ctx, key, strconv.FormatInt(compId, 10)).Err()
}
