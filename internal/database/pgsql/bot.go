package pgsql

import (
	"context"

	"github.com/botscubes/bot-service/internal/model"
)

func (db *Db) AddBot(m *model.Bot) (int64, error) {
	var id int64
	query := `INSERT INTO public.bot (user_id, token, title, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	if err := db.Pool.QueryRow(
		context.Background(), query, m.User_id, m.Token, m.Title, m.Status,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) CheckBotExist(user_id int64, bot_id int64) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE id = $1 AND user_id = $2) AS "exists";`
	if err := db.Pool.QueryRow(
		context.Background(), query, bot_id, user_id,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) CheckTokenExist(token *string) (bool, error) {
	var c bool
	query := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE token = $1) AS "exists";`
	if err := db.Pool.QueryRow(
		context.Background(), query, token,
	).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) GetBotToken(user_id int64, bot_id int64) (*string, error) {
	var data string
	query := `SELECT token FROM public.bot WHERE id = $1 AND user_id = $2;`
	if err := db.Pool.QueryRow(
		context.Background(), query, bot_id, user_id,
	).Scan(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (db *Db) SetBotToken(user_id int64, bot_id int64, token *string) error {
	query := `UPDATE public.bot SET token = $1 WHERE id = $2 AND user_id = $3;`
	_, err := db.Pool.Exec(context.Background(), query, token, bot_id, user_id)
	return err
}
