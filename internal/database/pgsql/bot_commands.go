package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddCommand(botId int64, m *model.Command) (int64, error) {
	var id int64
	query := `INSERT INTO ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			("type", "data", component_id, next_step_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Type, m.Data, m.ComponentId, m.NextStepId, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) SetNextStepCommand(botId int64, commandId int64, nextStepId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET next_step_id = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, nextStepId, commandId)
	return err
}

func (db *Db) CheckCommandExist(botId int64, compId int64, commandId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			WHERE id = $1 AND component_id = $2 AND status = $3) AS "exists";`

	if err := db.Pool.QueryRow(
		context.Background(), query, commandId, compId, model.StatusCommandActive,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) DelNextStepCommand(botId int64, commandId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET next_step_id = null WHERE id = $1;`

	_, err := db.Pool.Exec(context.Background(), query, commandId)
	return err
}

func (db *Db) DelCommandsByCompId(botId int64, compId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET status = $1 WHERE component_id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, model.StatusCommandDel, compId)
	return err
}

func (db *Db) DelCommand(botId int64, commandId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET status = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, model.StatusCommandDel, commandId)
	return err
}

func (db *Db) UpdCommand(botId int64, commandId int64, t *string, data *string) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET "type" = $1, "data" = $2 WHERE id = $3;`

	_, err := db.Pool.Exec(context.Background(), query, t, data, commandId)
	return err
}

func (db *Db) DelAllCommands(botId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET status = $1;`

	_, err := db.Pool.Exec(context.Background(), query, model.StatusCommandDel)
	return err
}
