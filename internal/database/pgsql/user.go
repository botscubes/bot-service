package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
)

// Statuses
var (
	StatusUserActive = 0
)

// TODO: check status

func (db *Db) AddUser(botId int64, m *model.User) (int64, error) {
	var id int64
	prefix := config.PrefixSchema + strconv.FormatInt(botId, 10)

	query := `INSERT INTO ` + prefix + `.user 
	(tg_id, first_name, last_name, username, step_id, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	if err := db.Pool.QueryRow(
		context.Background(), query, m.TgId, m.FirstName, m.LastName, m.Username, m.StepId, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) CheckUserExistByTgId(botId int64, TgId int64) (bool, error) {
	var c bool
	prefix := config.PrefixSchema + strconv.FormatInt(botId, 10)

	query := `SELECT EXISTS(SELECT 1 FROM ` + prefix + `.user WHERE tg_id = $1) AS "exists";`
	if err := db.Pool.QueryRow(
		context.Background(), query, TgId,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) UserByTgId(userId int64, botId int64) (*model.User, error) {
	prefix := config.PrefixSchema + strconv.FormatInt(botId, 10)

	query := `SELECT id, tg_id, first_name, last_name, username, status
			FROM ` + prefix + `.user WHERE tg_id = $1;`

	var r model.User
	if err := db.Pool.QueryRow(
		context.Background(), query, botId, userId,
	).Scan(&r.Id, &r.TgId, &r.FirstName, &r.LastName, &r.Username, &r.Status); err != nil {
		return nil, err
	}

	return &r, nil
}
