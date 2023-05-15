package pgsql

import (
	"context"

	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddBot(m *model.Bot) (int64, error) {
	var id int64
	query := `INSERT INTO public.bot (user_id, token, title, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	if err := db.Pool.QueryRow(
		context.Background(), query, m.UserId, m.Token, m.Title, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) CheckBotExist(userId int64, botId int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE id = $1 AND user_id = $2 AND status = $3) AS "exists";`
	if err := db.Pool.QueryRow(
		context.Background(), query, botId, userId, model.StatusBotActive,
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

	query := `SELECT id, title FROM public.bot WHERE status = $1 AND user_id = $2 ORDER BY id;`

	rows, err := db.Pool.Query(context.Background(), query, model.StatusBotActive, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.Bot
		if err = rows.Scan(&r.Id, &r.Title); err != nil {
			return nil, err
		}

		data = append(data, &r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &data, nil
}
