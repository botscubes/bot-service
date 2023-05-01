package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddComponent(botId int64, m *model.Component) (int64, error) {
	var id int64
	query := `INSERT INTO ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
			("data", keyboard, next_step_id, "position", status) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Data, m.Keyboard, m.NextStepId, m.Position, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// func (db *Db) GetComponent() {
// 	data := new(model.Component)

// 	query := `SELECT * FROM bot_41.component WHERE id = 5;`
// 	if err := db.Pool.QueryRow(
// 		context.Background(), query,
// 	).Scan(&data.Id, &data.Data, &data.Keyboard, &data.NextId, &data.Position, &data.Status); err != nil {
// 		log.Debug("err")
// 		log.Debug(err)
// 	}

// 	log.Debug(data)

// 	log.Debug("23")
// }

func (db *Db) AddCommand(botId int64, m *model.Command) (int64, error) {
	var id int64
	query := `INSERT INTO ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command
			("type", "data", component_id, next_step_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Type, m.Data, m.ComponentId, m.NextStepId, m.Status,
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

func (db *Db) SetNextStepForComponent(botId int64, compId int64, nextStepId int64) error {
	query := `UPDATE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET next_step_id = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, nextStepId, compId)
	return err
}

func (db *Db) CheckCommandExist(botId int64, compId int64, commandId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command
			WHERE id = $1 AND component_id = $2) AS "exists";`

	if err := db.Pool.QueryRow(
		context.Background(), query, commandId, compId,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) SetNextStepForCommand(botId int64, commandId int64, nextStepId int64) error {
	query := `UPDATE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET next_step_id = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, nextStepId, commandId)
	return err
}
