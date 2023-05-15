package redis

import (
	"context"
	"strconv"
)

func (rdb *Rdb) DelBotData(botId int64) {
	ctx := context.Background()
	key := "bot" + strconv.FormatInt(botId, 10) + "*"

	iter := rdb.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		_ = rdb.Del(ctx, iter.Val())
	}
}
