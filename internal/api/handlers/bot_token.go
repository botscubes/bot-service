package handlers

import (
	"strconv"

	"github.com/goccy/go-json"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/valyala/fasthttp"
)

type setTokenReq struct {
	Token *string `json:"token"`
}

func SetToken(db *pgsql.Db) fasthttp.RequestHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		var data setTokenReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: setToken] - Serialization error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetToken] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetToken] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

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

		existBot, err := db.CheckBotExist(userId, botId)
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

		oldToken, err := db.GetBotToken(userId, botId)
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

		if err = db.SetBotToken(userId, botId, token); err != nil {
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
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DeleteToken] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DeleteToken] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
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

		if err = db.SetBotToken(userId, botId, &token); err != nil {
			log.Error("[API: deleteToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}
