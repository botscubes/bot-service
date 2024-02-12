package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddConnection(botId int64, groupId int64, m *model.Connection) error {

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

	schema := prefixSchema + strconv.FormatInt(botId, 10)
	idx := "{\"" + strconv.FormatInt(*m.SourceComponentId, 10) + " " + *m.SourcePointName + "\"}"
	query := `
		UPDATE ` + schema + `.component 
		SET connection_points = JSONB_SET(connection_points, $1, $2) 
		WHERE group_id = $3 AND component_id = $4;`

	_, err = tx.Exec(
		ctx, query, idx, m.ConnectionPoint, groupId, m.TargetComponentId,
	)
	if err != nil {
		return err
	}

	idx = "{\"" + *m.SourcePointName + "\"}"
	query = `
		UPDATE ` + schema + `.component 
		SET outputs = JSONB_SET(outputs, $1, $2) 
		WHERE group_id = $3 AND component_id = $4;`

	_, err = tx.Exec(
		ctx, query, idx, m.TargetComponentId, groupId, m.SourceComponentId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) GetTargetComponentId(botId int64, groupId int64, m *model.SourceConnectionPoint) (*int64, error) {

	schema := prefixSchema + strconv.FormatInt(botId, 10)
	idx := *m.SourcePointName
	query := `
		SELECT outputs -> $1 FROM ` + schema + `.component 
		WHERE group_id = $2 AND component_id = $3;`

	var targetComponentId *int64
	if err := db.Pool.QueryRow(
		context.Background(), query, idx, groupId, m.SourceComponentId,
	).Scan(&targetComponentId); err != nil {
		return nil, err
	}
	return targetComponentId, nil
}

func (db *Db) DeleteConnection(botId int64, groupId int64, m *model.SourceConnectionPoint, targetComponentId int64) error {
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
	schema := prefixSchema + strconv.FormatInt(botId, 10)

	idx := "{\"" + strconv.FormatInt(*m.SourceComponentId, 10) + " " + *m.SourcePointName + "\"}"
	query := `
		UPDATE ` + schema + `.component 
		SET connection_points = JSONB_SET(connection_points, $1, 'null') 
		WHERE group_id = $2 AND component_id = $3;`

	_, err = tx.Exec(
		ctx, query, idx, groupId, targetComponentId,
	)
	if err != nil {
		return err
	}

	idx = "{\"" + *m.SourcePointName + "\"}"
	query = `
		UPDATE ` + schema + `.component 
		SET outputs = JSONB_SET(outputs, $1, 'null') 
		WHERE group_id = $2 AND component_id = $3;`

	_, err = tx.Exec(
		ctx, query, idx, groupId, m.SourceComponentId,
	)
	if err != nil {
		return err
	}
	return nil

}
