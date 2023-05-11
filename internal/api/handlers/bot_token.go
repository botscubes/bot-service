package handlers

import (
	"strconv"

	"github.com/goccy/go-json"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	fh "github.com/valyala/fasthttp"
)

type setBotTokenReq struct {
	Token *string `json:"token"`
}

func SetBotToken(db *pgsql.Db) reqHandler {
	// TODO: check bot is started
	return func(ctx *fh.RequestCtx) {
		var data setBotTokenReq

		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// check token
		token := data.Token
		if token == nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.MissingParam("token")))
			return
		}

		if !validateToken(*token) {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrIncorrectTokenFormat))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot token installed
		oldToken, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if oldToken != nil && *oldToken != "" {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrTokenAlreadyInstalled))
			return
		}

		// check token exists
		existToken, err := db.CheckBotTokenExist(token)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if existToken {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrTokenAlreadyExists))
			return
		}

		if err = db.SetBotToken(userId, botId, token); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func DeleteBotToken(db *pgsql.Db) reqHandler {
	// TODO: check bot is started
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		token := ""

		if err = db.SetBotToken(userId, botId, &token); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}
