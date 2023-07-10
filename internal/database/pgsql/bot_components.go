package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddComponent(botId int64, m *model.Component) (int64, error) {
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

func (db *Db) CheckComponentExist(botId int64, compId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			WHERE id = $1 AND status = $2) AS "exists";`

	if err := db.Pool.QueryRow(
		context.Background(), query, compId, model.StatusComponentActive,
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

func (db *Db) ComponentsForEd(botId int64) (*[]*model.Component, error) {
	var data []*model.Component

	query := `SELECT id, data, keyboard, ARRAY(
				SELECT jsonb_build_object('id', id, 'data', data, 'type', type, 'componentId', component_id, 'nextStepId', next_step_id)
				FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
				WHERE component_id = t.id AND status = $1 ORDER BY id
			), next_step_id, is_main, position, status
			FROM ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component t
			WHERE status = $2 ORDER BY id;`

	rows, err := db.Pool.Query(context.Background(), query, model.StatusCommandActive, model.StatusComponentActive)
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

func (db *Db) DelNextStepComponentByNS(botId int64, nextStep int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET next_step_id = null WHERE next_step_id = $1;`

	_, err := db.Pool.Exec(context.Background(), query, nextStep)
	return err
}

func (db *Db) DelComponent(botId int64, compId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET status = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, model.StatusComponentDel, compId)
	return err
}

func (db *Db) DelSetOfComponents(botId int64, data *[]int64) error {
	var err error
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

	// delete components
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET status = $1 WHERE id = ANY($2::bigint[]);`
	if _, err = tx.Exec(ctx, query, model.StatusComponentDel, data); err != nil {
		return err
	}

	// delete component commands
	query = `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET status = $1 WHERE component_id = ANY($2::bigint[]);`
	if _, err = tx.Exec(ctx, query, model.StatusCommandDel, data); err != nil {
		return err
	}

	// delete component next steps, that reference these components
	query = `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET next_step_id = null WHERE next_step_id = ANY($1::bigint[]);`
	if _, err = tx.Exec(ctx, query, data); err != nil {
		return err
	}

	// delete command next steps, that reference these components
	query = `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.command
			SET next_step_id = null WHERE next_step_id = ANY($1::bigint[]);`
	if _, err = tx.Exec(ctx, query, data); err != nil {
		return err
	}

	return err
}

func (db *Db) DelAllComponents(botId int64) error {
	query := `UPDATE ` + prefixSchema + strconv.FormatInt(botId, 10) + `.component
			SET status = $1 WHERE id != $2;`

	_, err := db.Pool.Exec(context.Background(), query, model.StatusComponentDel, config.MainComponentId)
	return err
}

func (db *Db) UpdComponentPosition(botId int64, compId int64, pos *model.Point) error {
	prefix := prefixSchema + strconv.FormatInt(botId, 10)
	query := `UPDATE ` + prefix + `.component SET "position" = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, pos, compId)
	return err
}

func (db *Db) UpdComponentData(botId int64, compId int64, data *model.ComponentData) error {
	prefix := prefixSchema + strconv.FormatInt(botId, 10)
	query := `UPDATE ` + prefix + `.component SET "data" = $1 WHERE id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, data, compId)
	return err
}
