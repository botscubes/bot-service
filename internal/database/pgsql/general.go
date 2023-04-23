package pgsql

import (
	"context"
	"strconv"
)

const (
	prefixSchema = "bot_"
)

func (db *Db) CreateSchema(bot_id int64) error {
	query := `CREATE SCHEMA IF NOT EXISTS ` + prefixSchema + strconv.FormatInt(bot_id, 10)
	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *Db) CreateBotUserTable(bot_id int64) error {
	query := `CREATE TABLE ` + prefixSchema + strconv.FormatInt(bot_id, 10) + `.user
	(
		id bigserial NOT NULL,
		tg_id bigint NOT NULL,
		first_name text,
		last_name text,
		username text,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *Db) CreateBotStructureTable(bot_id int64) error {
	query := `CREATE TABLE ` + prefixSchema + strconv.FormatInt(bot_id, 10) + `.structure
	(
		id bigserial NOT NULL,
		component jsonb NOT NULL,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}
