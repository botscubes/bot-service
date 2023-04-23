package handlers

import (
	"encoding/json"
	"regexp"

	"github.com/botscubes/bot-service/internal/api/errors"
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

func setToken(db *pgsql.Db) fasthttp.RequestHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data setTokenReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: setToken] - Serialisation error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: setToken] - (bot_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: setToken] - (user_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		token := data.Token
		if token == nil {
			log.Debug("[API: setToken] - token is misssing")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if !validateToken(*token) {
			log.Debug("[API: setToken] - Incorrect token")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrIncorrectTokenFormat))
			return
		}

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: setToken] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if !existBot {
			log.Debug("[API: setToken] - bot not found")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		existToken, err := db.CheckTokenExist(token)
		if err != nil {
			log.Debug("[API: setToken] - [db: CheckTokenExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if existToken {
			log.Debug("[API: setToken] - token exists")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrTokenAlreadyExists))
			return
		}

		oldToken, err := db.GetBotToken(bot_id)
		if err != nil {
			log.Debug("[API: setToken] - [db: GetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if oldToken != nil && *oldToken != "" {
			log.Debug("[API: setToken] - Token is already installed")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrTokenAlreadyInstalled))
			return
		}

		if err = db.SetBotToken(bot_id, token); err != nil {
			log.Debug("[API: setToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		dataRes := &messageRes{
			Message: "Token installed",
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}

func deleteToken(db *pgsql.Db) fasthttp.RequestHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data deleteTokenReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: deleteToken] - Serialisation error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: deleteToken] - (bot_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: deleteToken] - (user_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: deleteToken] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if !existBot {
			log.Debug("[API: deleteToken] - bot not found")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		token := ""

		if err = db.SetBotToken(bot_id, &token); err != nil {
			log.Debug("[API: deleteToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		dataRes := &messageRes{
			Message: "Token deleted",
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}
