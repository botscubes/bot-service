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

type setBotTokenReq struct {
	Token *string `json:"token"`
}

func SetBotToken(db *pgsql.Db) reqHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		var data setBotTokenReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: SetBotToken] - Serialization error;", err)
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
			log.Debug("[API: SetBotToken] - token is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if !validateToken(*token) {
			log.Debug("[API: SetBotToken] - Incorrect token")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrIncorrectTokenFormat))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: SetBotToken] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: SetBotToken] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		oldToken, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Debug("[API: SetBotToken] - [db: GetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if oldToken != nil && *oldToken != "" {
			log.Debug("[API: SetBotToken] - Token is already installed")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrTokenAlreadyInstalled))
			return
		}

		existToken, err := db.CheckBotTokenExist(token)
		if err != nil {
			log.Debug("[API: SetBotToken] - [db: CheckBotTokenExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if existToken {
			log.Debug("[API: SetBotToken] - token exists")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrTokenAlreadyExists))
			return
		}

		if err = db.SetBotToken(userId, botId, token); err != nil {
			log.Error("[API: SetBotToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}

func DeleteBotToken(db *pgsql.Db) reqHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DeleteBotToken] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DeleteBotToken] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: DeleteBotToken] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: DeleteBotToken] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		token := ""

		if err = db.SetBotToken(userId, botId, &token); err != nil {
			log.Error("[API: DeleteBotToken] - [db: SetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}
