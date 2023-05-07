package redis

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
	"github.com/redis/go-redis/v9"
)

func (rdb *Rdb) SetUser(botId int64, user *model.User) error {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":user"

	return rdb.HSet(ctx, key, strconv.FormatInt(user.TgId, 10), user).Err()
}

func (rdb *Rdb) GetUser(botId int64, userID int64) (*model.User, error) {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":user"

	user := &model.User{}

	result, err := rdb.HGet(ctx, key, strconv.FormatInt(userID, 10)).Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	if result == "" {
		return nil, nil
	}

	if err := user.UnmarshalBinary([]byte(result)); err != nil {
		return nil, err
	}

	return user, nil
}
