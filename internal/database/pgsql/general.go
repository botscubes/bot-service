package pgsql

import (
	"context"
	"strconv"
)

func (db *Db) CreateBotSchema(botId int64) error {
	query := `CALL create_bot_schema(` + strconv.FormatInt(botId, 10) + `);`
	_, err := db.Pool.Exec(context.Background(), query)
	return err
}
