package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
)

func (db *Db) NewComponent(bot_id int64, m *model.Component) (int64, error) {
	var id int64
	query := `INSERT INTO` + config.PrefixSchema + strconv.FormatInt(bot_id, 10) + `.structure
			("data", keyboard, next_id, "position", status) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Data, m.Keyboard, m.NextId, m.Position, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) GetComponent() {
	data := new(model.Component)

	query := `SELECT * FROM bot_41.structure WHERE id = 5;`
	if err := db.Pool.QueryRow(
		context.Background(), query,
	).Scan(&data.Id, &data.Data, &data.Keyboard, &data.NextId, &data.Position, &data.Status); err != nil {
		log.Debug("err")
		log.Debug(err)
	}

	log.Debug(data)

	log.Debug(123)
}
