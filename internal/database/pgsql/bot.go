package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) CreateBot(m *model.Bot, mc *model.Component) (botId int64, componentId int64, err error) {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// create bot
	query := `INSERT INTO public.bot (user_id, token, title, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	if err = tx.QueryRow(
		ctx, query, m.UserId, m.Token, m.Title, m.Status,
	).Scan(&botId); err != nil {
		return 0, 0, err
	}

	// create schema for bot
	query = `CALL create_bot_schema(` + strconv.FormatInt(botId, 10) + `);`
	if _, err = tx.Exec(ctx, query); err != nil {
		return 0, 0, err
	}
	var groupId int
	bot := prefixSchema + strconv.FormatInt(botId, 10)
	query = `INSERT INTO ` + bot + `.component_group
			(name) VALUES ('main') RETURNING id;`
	if err = tx.QueryRow(
		ctx, query,
	).Scan(&groupId); err != nil {
		return 0, 0, err
	}

	// add component
	query = `INSERT INTO ` + bot + `.component
			(type, next_id, path, position, group_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	if err = tx.QueryRow(
		ctx, query, mc.Type, mc.NextComponentId, mc.Path, mc.Position, groupId,
	).Scan(&componentId); err != nil {
		return 0, 0, err
	}

	return botId, componentId, nil
}

func (db *Db) DeleteBot(userId int64, botId int64) error {
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
	query := `DELETE FROM public.bot WHERE id = $1 AND user_id = $2;`
	_, err = db.Pool.Exec(ctx, query, botId, userId)
	if err != nil {
		return err
	}
	query = `DROP SCHEMA IF EXISTS bot_` + strconv.FormatInt(botId, 10) + ` CASCADE;`
	_, err = db.Pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (db *Db) CheckBotExist(userId int64, botId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE id = $1 AND user_id = $2 AND (status = $3 OR status = $4)) AS "exists";`
	if err := db.Pool.QueryRow(
		context.Background(), query, botId, userId, model.StatusBotStopped, model.StatusBotRunning,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) CheckBotTokenExist(token *string) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE token = $1) AS "exists";`
	if err := db.Pool.QueryRow(
		context.Background(), query, token,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) GetBotToken(userId int64, botId int64) (*string, error) {
	var data string
	query := `SELECT token FROM public.bot WHERE id = $1 AND user_id = $2;`
	if err := db.Pool.QueryRow(
		context.Background(), query, botId, userId,
	).Scan(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (db *Db) SetBotToken(userId int64, botId int64, token *string) error {
	query := `UPDATE public.bot SET token = $1 WHERE id = $2 AND user_id = $3;`
	_, err := db.Pool.Exec(context.Background(), query, token, botId, userId)
	return err
}

func (db *Db) UserBots(userId int64) (*[]*model.Bot, error) {
	data := []*model.Bot{}

	query := `SELECT id, title, status FROM public.bot WHERE user_id = $1 AND (status = $2 OR status = $3) ORDER BY id;`

	rows, err := db.Pool.Query(context.Background(), query, userId, model.StatusBotStopped, model.StatusBotRunning)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.Bot
		if err = rows.Scan(&r.Id, &r.Title, &r.Status); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}

func (db *Db) GetBotStatus(botId int64, userId int64) (model.BotStatus, error) {
	var data model.BotStatus
	query := `SELECT status FROM public.bot WHERE id = $1 AND user_id = $2;`
	if err := db.Pool.QueryRow(
		context.Background(), query, botId, userId,
	).Scan(&data); err != nil {
		return 0, err
	}

	return data, nil
}

func (db *Db) SetBotStatus(botId int64, userId int64, status model.BotStatus) error {
	query := `UPDATE public.bot SET status = $1 WHERE id = $2 AND user_id = $3;`
	_, err := db.Pool.Exec(context.Background(), query, status, botId, userId)
	return err
}
