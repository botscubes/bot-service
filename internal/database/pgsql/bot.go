package pgsql

import (
	"context"
)

func (db *Db) AddBot(user_id int64, token *string, title *string, status int) (int64, error) {
	var id int64
	queryInsert := `INSERT INTO public.bot (user_id, token, title, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	if err := db.Pool.QueryRow(context.Background(), queryInsert, user_id, token, title, status).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Db) CheckBotExist(user_id int64, bot_id int64) (bool, error) {
	var c bool
	queryInsert := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE id = $1 AND user_id = $2) AS "exists";`
	if err := db.Pool.QueryRow(context.Background(), queryInsert, bot_id, user_id).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) CheckTokenExist(token *string) (bool, error) {
	var c bool
	queryInsert := `SELECT EXISTS(SELECT 1 FROM public.bot WHERE token = $1) AS "exists";`
	if err := db.Pool.QueryRow(context.Background(), queryInsert, token).Scan(&c); err != nil {
		return false, err
	}

	return c, nil
}

func (db *Db) GetBotToken(bot_id int64) (*string, error) {
	var data string
	queryInsert := `SELECT token FROM public.bot WHERE id = $1;`
	if err := db.Pool.QueryRow(context.Background(), queryInsert, bot_id).Scan(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (db *Db) SetBotToken(bot_id int64, token *string) error {
	queryInsert := `UPDATE public.bot SET token = $1 WHERE id = $2;`
	_, err := db.Pool.Exec(context.Background(), queryInsert, token, bot_id)
	return err
}
