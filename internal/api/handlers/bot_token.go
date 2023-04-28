package handlers

import (
	"encoding/json"
	"regexp"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/api/schema"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/valyala/fasthttp"
)

const (
	tokenRegexp = `^\d{9,10}:[\w-]{35}$` //nolint:gosec
)

func validateToken(token string) bool {
	reg := regexp.MustCompile(tokenRegexp)
	return reg.MatchString(token)
}

func SetToken(db *pgsql.Db) fasthttp.RequestHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data schema.SetTokenReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: setToken] - Serialisation error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if data.BotId == nil {
			log.Debug("[API: setToken] bot_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: setToken] - (bot_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id := ctx.UserValue("user_id").(int64)

		token := data.Token
		if token == nil {
			log.Debug("[API: setToken] - token is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if !validateToken(*token) {
			log.Debug("[API: setToken] - Incorrect token")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrIncorrectTokenFormat))
			return
		}

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: setToken] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: setToken] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		oldToken, err := db.GetBotToken(user_id, bot_id)
		if err != nil {
			log.Debug("[API: setToken] - [db: GetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if oldToken != nil && *oldToken != "" {
			log.Debug("[API: setToken] - Token is already installed")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrTokenAlreadyInstalled))
			return
		}

		existToken, err := db.CheckTokenExist(token)
		if err != nil {
			log.Debug("[API: setToken] - [db: CheckTokenExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if existToken {
			log.Debug("[API: setToken] - token exists")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrTokenAlreadyExists))
			return
		}

		if err = db.SetBotToken(user_id, bot_id, token); err != nil {
			log.Error("[API: setToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}

func DeleteToken(db *pgsql.Db) fasthttp.RequestHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data schema.BotIdReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: deleteToken] - Serialisation error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if data.BotId == nil {
			log.Debug("[API: deleteToken] bot_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: deleteToken] - (bot_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id := ctx.UserValue("user_id").(int64)

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: deleteToken] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: deleteToken] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		token := ""

		if err = db.SetBotToken(user_id, bot_id, &token); err != nil {
			log.Error("[API: deleteToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}
