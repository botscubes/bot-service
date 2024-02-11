package pgsql

import (
	"context"
	"strconv"

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

	//jsonData, err := json.Marshal(m.ConnectionPoint)
	//if err != nil {
	//	return err
	//}
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

	return nil
}

func (db *Db) DeleteConnection(botId int64, groupId int64, m *model.DelConnectionReq) error {
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
