package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
)

//				id bigserial NOT NULL,
//	           type VARCHAR(20),
//	           component_id BIGINT NOT NULL UNIQUE,
//	           next_id BIGINT,
//	           path text NOT NULL DEFAULT '''',
//	           position POINT,
//	           group_id BIGINT,
//	           PRIMARY KEY (id),
//	           FOREIGN KEY(group_id)
func (db *Db) AddComponent(botId int64, groupId int64, m *model.Component) (int64, error) {

	schema := prefixSchema + strconv.FormatInt(botId, 10)
	var id int64
	query := `INSERT INTO ` + schema + `.component
			(type, path, position, group_id) VALUES ($1, $2, $3, $4) RETURNING component_id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.Type, m.Path, m.Position, groupId,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) DeleteComponent(botId int64, groupId int64, componentId int64) error {

	schema := prefixSchema + strconv.FormatInt(botId, 10)

	query := `DELETE FROM` + schema + `.component
			WHERE group_id = $1 AND component_id = $2;`

	_, err := db.Pool.Exec(context.Background(), query, groupId, componentId)
	return err

}
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

}

func (db *Db) CheckComponentExist(botId int64, groupId int64, compId int64) (bool, error) {
	schema := prefixSchema + strconv.FormatInt(botId, 10)
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + schema + `.component
			WHERE group_id = $1 AND component_id = $2);`

	if err := db.Pool.QueryRow(
		context.Background(), query, groupId, compId,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}
