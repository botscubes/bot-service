package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
)

func (db *Db) AddComponent(bot_id int64, m *model.Component) (int64, error) {
	var id int64
	query := `INSERT INTO ` + config.PrefixSchema + strconv.FormatInt(bot_id, 10) + `.component
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

	query := `SELECT * FROM bot_41.component WHERE id = 5;`
	if err := db.Pool.QueryRow(
		context.Background(), query,
	).Scan(&data.Id, &data.Data, &data.Keyboard, &data.NextId, &data.Position, &data.Status); err != nil {
		log.Debug("err")
		log.Debug(err)
	}

	log.Debug(data)

	log.Debug(123)
}

func (db *Db) AddCommand(bot_id int64, m *model.Command) (int64, error) {
	var id int64
	query := `INSERT INTO ` + config.PrefixSchema + strconv.FormatInt(bot_id, 10) + `.command
			("type", "data", next_id) VALUES ($1, $2, $3) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Type, m.Data, m.NextId,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) CheckComponentExist(botId int64, compId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
			WHERE id = $1) AS "exists";`

	if err := db.Pool.QueryRow(
		context.Background(), query, compId,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) SetNextIdForComponent(botId int64, compId int64, nextId int64) error {
	query := `UPDATE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET next_id = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, nextId, compId)
	return err
}
