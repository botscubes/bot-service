package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddBotComponent(botId int64, m *model.Component) (int64, error) {
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

func (db *Db) AddBotCommand(botId int64, m *model.Command) (int64, error) {
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

func (db *Db) CheckBotComponentExist(botId int64, compId int64) (bool, error) {
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

func (db *Db) CheckBotCommandExist(botId int64, compId int64, commandId int64) (bool, error) {
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

func (db *Db) GetBotComponents(botId int64) (*[]*model.Component, error) {
	var data []*model.Component
	status := 0

	query := `SELECT * FROM ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
			WHERE status = $1;`

	rows, err := db.Pool.Query(context.Background(), query, status)
	if err != nil {
		return nil, err
	}

	// WARN: status not used
	for rows.Next() {
		var r model.Component
		if err = rows.Scan(&r.Id, &r.Data, &r.Keyboard, &r.NextStepId, &r.Position, &r.Status); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}

func (db *Db) GetBotFullComponents(botId int64) (*[]*model.ComponentFull, error) {
	var data []*model.ComponentFull
	status := 0

	query := `SELECT id, data, keyboard, ARRAY(
				SELECT jsonb_build_object('id', id, 'data', data, 'type', type, 'component_id', component_id, 'next_step_id', next_step_id)
				FROM ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command where component_id = t.id
			), next_step_id, position, status FROM ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component t WHERE status = $1;`

	rows, err := db.Pool.Query(context.Background(), query, status)
	if err != nil {
		return nil, err
	}

	// WARN: status not used
	for rows.Next() {
		var r model.ComponentFull
		if err = rows.Scan(&r.Id, &r.Data, &r.Keyboard, &r.Commands, &r.NextStepId, &r.Position, &r.Status); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}

func (db *Db) GetBotCommands(botId int64) (*[]*model.Command, error) {
	var data []*model.Command
	status := 0

	query := `SELECT * FROM ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command
			WHERE status = $1;`

	rows, err := db.Pool.Query(context.Background(), query, status)
	if err != nil {
		return nil, err
	}

	// WARN: status not used
	for rows.Next() {
		var r model.Command
		if err = rows.Scan(&r.Id, &r.Type, &r.Data, &r.ComponentId, &r.NextStepId, &r.Status); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}
