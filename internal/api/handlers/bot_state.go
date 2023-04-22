package handlers

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
)

const (
	maxTitleLen = 50
)

func newBot(db *pgsql.Db) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data newBotReq
		err = json.Unmarshal(ctx.PostBody(), &data)
		if err != nil {
			log.Debug("[API: newBot] - Serialisation error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if data.UserId == nil {
			log.Debug("[API: newBot] user_id is misssing")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: newBot] - (user_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		title := data.Title
		token := ""
		status := 0

		if title == nil || *title == "" {
			log.Debug("[API: newBot] - title is misssing")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if utf8.RuneCountInString(*title) > maxTitleLen {
			log.Debug("[API: newBot] - title len > ", maxTitleLen)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidTitleLength))
			return
		}

		botId, err := db.AddBot(user_id, &token, title, status)
		if err != nil {
			log.Debug("[API: newBot] - [db: AddBot] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		log.Info(botId)

		dataRes := &newBotRes{
			Id: botId,
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}

func startBot(db *pgsql.Db, bots *map[string]*bot.TBot, server *telego.MultiBotWebhookServer, conf *config.BotConfig) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data startBotReq
		err = json.Unmarshal(ctx.PostBody(), &data)
		if err != nil {
			log.Errorf("[API: startBot] - Serialisation error;\n %s", err)
			doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrInvalidRequest)
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: startBot] - (bot_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: startBot] - (user_id) json.Number convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: startBot] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if !existBot {
			log.Debug("[API: startBot] - bot not found")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		token, err := db.GetBotToken(bot_id)
		if err != nil {
			log.Debug("[API: startBot] - [db: GetBotToken] error;", err)
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if token == nil || *token == "" {
			log.Debug("[API: startBot] - Token not found")
			doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrTokenNotFound))
			return
		}

		// TODO: Own token health check to get a specific error
		nbot, err := bot.NewBot(token)
		if err != nil {
			log.Debug("[API: startBot] ", err)
			doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrInvalidToken)
			return
		}

		if _, ok := (*bots)[*token]; ok {
			log.Debug("[API: startBot] Token already exist in bots map")
			doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrTokenExistInSystem)
			return
		}

		(*bots)[*token] = new(bot.TBot)
		(*bots)[*token].Bot = nbot

		err = (*bots)[*token].StartBot(conf.WebhookBase, conf.ListenAddress, server)
		if err != nil {
			log.Debug("[API: startBot] Start bot error ", err)
			doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrStartBot)
			return
		}

		dataRes := &messageRes{
			Message: "Bot started",
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}
