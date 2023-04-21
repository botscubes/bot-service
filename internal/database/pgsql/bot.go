package pgsql

import (
	"context"
	"fmt"
	"os"
)

func (db *Db) GetTest() {
	var version string
	err := db.Pool.QueryRow(context.Background(), "select version()").Scan(&version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(version)
}

func (db *Db) AddBot(user_id int64, token *string, title *string, status int) (int64, error) {
	var id int64
	queryInsert := `INSERT INTO public.bot (user_id, token, title, status) VALUES ($1, $2, $3, $4) RETURNING id;`
	err := db.Pool.QueryRow(context.Background(), queryInsert, user_id, token, title, status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
