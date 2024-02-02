package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
)

//			id bigserial NOT NULL,
//            type VARCHAR(20),
//            component_id BIGINT NOT NULL UNIQUE,
//            next_id BIGINT,
//            path text NOT NULL DEFAULT '''',
//            position POINT,
//            group_id BIGINT,
//            PRIMARY KEY (id),
//            FOREIGN KEY(group_id)

func (db *Db) GetComponents(botId int64, groupId int64) ([]*model.Component, error) {

	schema := prefixSchema + strconv.FormatInt(botId, 10)
	query := `
		SELECT 
			
			component_id, 
			type, 
			next_id, 
			path, 
			position 
		FROM ` + schema + `.component WHERE group_id = $1;`

	rows, err := db.Pool.Query(context.Background(), query, groupId)
	if err != nil {
		return nil, err
	}

	var data []*model.Component
	// WARN: status not used
	for rows.Next() {
		var c model.Component

		if err = rows.Scan(&c.Id, &c.Type, &c.NextComponentId, &c.Path, &c.Position); err != nil {
			return nil, err
		}

		data = append(data, &c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return data, nil
	// query := `INSERT INTO ` + +`.component
	//
	//	("data", keyboard, next_step_id, is_main,"position", status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	//
	// if err := tx.QueryRow(
	//
	//	ctx, query, m.Data, m.Keyboard, m.NextStepId, m.IsMain, m.Position, m.Status,
	//
	//	).Scan(&id); err != nil {
	//		return 0, err
	//	}
	//
	// return id, nil
}
