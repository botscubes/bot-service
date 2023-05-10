package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
)

// Statuses
var (
	StatusComponentActive = 0
	StatusComponentDel    = 1
)

var (
	StatusCommandActive = 0
	StatusCommandDel    = 1
)

func (db *Db) AddBotComponent(botId int64, m *model.Component) (int64, error) {
	var id int64
	query := `INSERT INTO ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			("data", keyboard, next_step_id, is_main,"position", status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Data, m.Keyboard, m.NextStepId, m.IsMain, m.Position, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

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

func (db *Db) CheckComponentExist(botId int64, compId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			WHERE id = $1 AND status = $2) AS "exists";`

	if err := db.Pool.QueryRow(
		context.Background(), query, compId, StatusComponentActive,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) SetNextStepComponent(botId int64, compId int64, nextStepId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET next_step_id = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, nextStepId, compId)
	return err
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
		context.Background(), query, commandId, compId, StatusCommandActive,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) BotComponentsForEd(botId int64) (*[]*model.Component, error) {
	var data []*model.Component

	query := `SELECT id, data, keyboard, ARRAY(
				SELECT jsonb_build_object('id', id, 'data', data, 'type', type, 'componentId', component_id, 'nextStepId', next_step_id)
				FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
				WHERE component_id = t.id AND status = $1 ORDER BY id
			), next_step_id, is_main, position, status
			FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component t
			WHERE status = $2 ORDER BY id;`

	rows, err := db.Pool.Query(context.Background(), query, StatusCommandActive, StatusComponentActive)
	if err != nil {
		return nil, err
	}

	// WARN: status not used
	for rows.Next() {
		var r model.Component
		r.Commands = &model.Commands{}

		if err = rows.Scan(&r.Id, &r.Data, &r.Keyboard, r.Commands, &r.NextStepId, &r.IsMain, &r.Position, &r.Status); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}

func (db *Db) DelNextStepComponent(botId int64, compId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET next_step_id = null WHERE id = $1;`

	_, err := db.Pool.Exec(context.Background(), query, compId)
	return err
}

func (db *Db) DelNextStepCommand(botId int64, commandId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET next_step_id = null WHERE id = $1;`

	_, err := db.Pool.Exec(context.Background(), query, commandId)
	return err
}

func (db *Db) DelBotComponent(botId int64, compId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET status = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, StatusComponentDel, compId)
	return err
}

func (db *Db) DelCommandsByCompId(botId int64, compId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET status = $1 WHERE component_id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, StatusCommandDel, compId)
	return err
}

func (db *Db) DelCommand(botId int64, commandId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET status = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, StatusCommandDel, commandId)
	return err
}

func (db *Db) BotComponentsForBot(botId int64) (*[]*model.Component, error) {
	// TODO: REMOVE POSITION !
	// Change return value to *map[int64]*model.Component???

	var data []*model.Component

	query := `SELECT id, data, keyboard, ARRAY(
			SELECT jsonb_build_object('id', id, 'data', data, 'type', type, 'componentId', component_id, 'nextStepId', next_step_id)
				FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
				WHERE component_id = t.id AND status = $1 ORDER BY id
			), next_step_id, is_main, position, status
			FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component t
			WHERE status = $2 ORDER BY id;`

	rows, err := db.Pool.Query(context.Background(), query, StatusCommandActive, StatusComponentActive)
	if err != nil {
		return nil, err
	}

	// WARN: status not used
	for rows.Next() {
		var r model.Component
		r.Commands = &model.Commands{}
		if err = rows.Scan(&r.Id, &r.Data, &r.Keyboard, r.Commands, &r.NextStepId, &r.IsMain, &r.Position, &r.Status); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}

func (db *Db) ComponentForBot(botId int64, compID int64) (*model.Component, error) {
	// TODO: REMOVE POSITION !

	prefix := prefixSchema + strconv.FormatInt(botId, 10)

	query := `SELECT id, data, keyboard, ARRAY(
		SELECT jsonb_build_object('id', id, 'data', data, 'type', type, 'componentId', component_id, 'nextStepId', next_step_id)
		FROM ` + prefix + `.command
		WHERE component_id = t.id AND status = $1 ORDER BY id
	), next_step_id, is_main, position, status
	FROM ` + prefix + `.component t
	WHERE status = $2 AND id = $3 ORDER BY id;`

	// WARN: status not used
	var r model.Component
	r.Commands = &model.Commands{}
	if err := db.Pool.QueryRow(
		context.Background(), query, StatusCommandActive, StatusComponentActive, compID,
	).Scan(&r.Id, &r.Data, &r.Keyboard, r.Commands, &r.NextStepId, &r.IsMain, &r.Position, &r.Status); err != nil {
		return nil, err
	}

	return &r, nil
}
