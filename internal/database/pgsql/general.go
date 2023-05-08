package pgsql

import (
	"context"
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/pkg/log"
)

func (db *Db) CreateBotSchema(botId int64) error {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			if err.Error() != "tx is closed" {
				log.Error(err)
			}
		}
	}()

	newSchemaQ := `CREATE SCHEMA IF NOT EXISTS ` + config.PrefixSchema + strconv.FormatInt(botId, 10)

	_, err = tx.Exec(ctx, newSchemaQ)
	if err != nil {
		return err
	}

	userTableQ := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.user
	(
		id bigserial NOT NULL,
		tg_id bigint NOT NULL,
		first_name text,
		last_name text,
		username text,
		step_id bigint NOT NULL DEFAULT 1,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	_, err = tx.Exec(ctx, userTableQ)
	if err != nil {
		return err
	}

	componentQ := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.component
	(
		id bigserial NOT NULL,
		data jsonb,
		keyboard jsonb,
		next_step_id bigint,
		is_main boolean DEFAULT false,
		position point,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	_, err = tx.Exec(ctx, componentQ)
	if err != nil {
		return err
	}

	commandQ := `CREATE TABLE ` + config.PrefixSchema + strconv.FormatInt(botId, 10) + `.command
	(
		id bigserial NOT NULL,
		type character varying(20) NOT NULL,
		data text NOT NULL,
		component_id bigint NOT NULL,
		next_step_id bigint,
		status integer NOT NULL,
		PRIMARY KEY (id)
	)`

	_, err = tx.Exec(ctx, commandQ)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
