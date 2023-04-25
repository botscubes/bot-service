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
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: newBot] - Serialisation error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if data.UserId == nil {
			log.Debug("[API: newBot] user_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: newBot] - (user_id) json.Number convertation to int64 error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		title := data.Title
		token := ""
		status := 0

		if title == nil || *title == "" {
			log.Debug("[API: newBot] - title is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if utf8.RuneCountInString(*title) > maxTitleLen {
			log.Debug("[API: newBot] - title len > ", maxTitleLen)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidTitleLength))
			return
		}

		// TODO: Mb combine into one query (for rollback all on error)

		botId, err := db.AddBot(user_id, &token, title, status)
		if err != nil {
			log.Error("[API: newBot] - [db: AddBot] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if err := db.CreateSchema(botId); err != nil {
			log.Error("[API: newBot] - [db: CreateSchema] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if err := db.CreateBotUserTable(botId); err != nil {
			log.Error("[API: newBot] - [db: CreateBotUserTable] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if err := db.CreateBotStructureTable(botId); err != nil {
			log.Error("[API: newBot] - [db: CreateBotStructureTable] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		dataRes := &newBotRes{
			Id: botId,
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}

func startBot(db *pgsql.Db, bots *map[string]*bot.TBot, server *telego.MultiBotWebhookServer, conf *config.BotConfig) fasthttp.RequestHandler {
	// check bot already started
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data startBotReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Error("[API: startBot] - Serialisation error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, &errors.ErrInvalidRequest)
			return
		}

		if data.BotId == nil {
			log.Debug("[API: startBot] bot_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if data.UserId == nil {
			log.Debug("[API: startBot] user_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: startBot] - (bot_id) json.Number convertation to int64 error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: startBot] - (user_id) json.Number convertation to int64 error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: startBot] - [db: CheckBotExist] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: startBot] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		token, err := db.GetBotToken(bot_id)
		if err != nil {
			log.Debug("[API: startBot] - [db: GetBotToken] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if token == nil || *token == "" {
			log.Debug("[API: startBot] - Token not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrTokenNotFound))
			return
		}

		if _, ok := (*bots)[*token]; !ok {
			// TODO: Own token health check to get a specific error
			nbot, err := bot.NewBot(token)
			if err != nil {
				log.Debug("[API: startBot] ", err)
				doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrInvalidToken)
				return
			}

			(*bots)[*token] = new(bot.TBot)
			(*bots)[*token].Bot = nbot
		}

		if err = (*bots)[*token].StartBot(conf.WebhookBase, conf.ListenAddress, server); err != nil {
			log.Debug("[API: startBot] Start bot error ", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, &errors.ErrStartBot)
			return
		}

		dataRes := &messageRes{
			Message: "Bot started",
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}

func stopBot(db *pgsql.Db, bots *map[string]*bot.TBot) fasthttp.RequestHandler {
	// TODO: check bot is running
	// check bot already stopped
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data stopBotReq
		if err = json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Error("[API: stopBot] - Serialisation error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, &errors.ErrInvalidRequest)
			return
		}

		if data.BotId == nil {
			log.Debug("[API: stopBot] bot_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		if data.UserId == nil {
			log.Debug("[API: stopBot] user_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		bot_id, err := data.BotId.Int64()
		if err != nil {
			log.Debug("[API: stopBot] - (bot_id) json.Number convertation to int64 error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		user_id, err := data.UserId.Int64()
		if err != nil {
			log.Debug("[API: stopBot] - (user_id) json.Number convertation to int64 error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(user_id, bot_id)
		if err != nil {
			log.Debug("[API: stopBot] - [db: CheckBotExist] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: stopBot] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		token, err := db.GetBotToken(bot_id)
		if err != nil {
			log.Debug("[API: stopBot] - [db: GetBotToken] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if token == nil || *token == "" {
			log.Debug("[API: stopBot] - Token not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrTokenNotFound))
			return
		}

		if _, ok := (*bots)[*token]; !ok {
			log.Debug("[API: startBot] Bot not found in bots map")
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrBotNotFoundInSystem))
			return
		}

		if err := (*bots)[*token].StopBot(false); err != nil {
			log.Error("[API: stopBot] - Bot stop:\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrStopBot))
			return
		}

		dataRes := &messageRes{
			Message: "Bot stopped",
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}
