package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
)

func (db *Db) CreateSchema(botId int64) error {
	query := `CREATE SCHEMA IF NOT EXISTS ` + config.PrefixSchema + strconv.FormatInt(botId, 10)
	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *Db) CreateBotUserTable(botId int64) error {
	query := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.user
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

func (db *Db) CreateBotComponentTable(botId int64) error {
	query := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
	(
		id bigserial NOT NULL,
		data jsonb,
		keyboard jsonb,
		next_step_id bigint,
		is_start boolean DEFAULT false,
		position point,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *Db) CreateBotCommandTable(botId int64) error {
	query := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command
	(
		id bigserial NOT NULL,
		type character varying(20) NOT NULL,
		data text NOT NULL,
		component_id bigint NOT NULL,
		next_step_id bigint,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	if _, err := db.Pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}
