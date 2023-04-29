package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
)

func (db *Db) CreateSchema(bot_id int64) error {
	query := `CREATE SCHEMA IF NOT EXISTS ` + config.PrefixSchema + strconv.FormatInt(bot_id, 10)
	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *Db) CreateBotUserTable(bot_id int64) error {
	query := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(bot_id, 10) + `.user
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
	query := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(bot_id, 10) + `.structure
	(
		id bigserial NOT NULL,
		data jsonb,
		keyboard jsonb,
		next_id bigint,
		position point,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *Db) CreateBotCommandTable(bot_id int64) error {
	query := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(bot_id, 10) + `.command
	(
		id bigserial NOT NULL,
		type character varying(20) NOT NULL,
		data text NOT NULL,
		next_id bigint,
		PRIMARY KEY (id)
	)`

	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}
