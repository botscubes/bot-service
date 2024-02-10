package pgsql

import (
	"context"

	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddConnection(botId int64, groupId int64, m *model.Connection) error {
	// Проверить, существует ли соединение, если да, то нужно заменить, иначе добавить

	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	return nil
}
