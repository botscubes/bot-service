package redis

import (
	"context"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func (rdb *Rdb) SetUserStep(botId int64, userID int64, stepID int64) error {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":user:" + strconv.FormatInt(userID, 10) + ":step"

	return rdb.Set(ctx, key, stepID, 0).Err()
}

func (rdb *Rdb) GetUserStep(botId int64, userID int64) (int64, error) {
	ctx := context.Background()

	key := "bot" + strconv.FormatInt(botId, 10) + ":user:" + strconv.FormatInt(userID, 10) + ":step"

	val, err := rdb.Get(ctx, key).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	if errors.Is(err, redis.Nil) {
		// todo: change errors
		return 0, errors.New("not found")
	}

	stepID, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return int64(stepID), nil
}
