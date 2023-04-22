package handlers

import (
	"encoding/json"
	"regexp"
	"unicode/utf8"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	fastRouter "github.com/fasthttp/router"
	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
)

// TODO: check user access
// Add check all id's for > 0

const (
	tokenRegexp = `^\d{9,10}:[\w-]{35}$` //nolint:gosec
	maxTitleLen = 50
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func validateToken(token string) bool {
	reg := regexp.MustCompile(tokenRegexp)
	return reg.MatchString(token)
}

func doJsonRes(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

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

func setToken(db *pgsql.Db) fasthttp.RequestHandler {
	// TODO: check bot is started
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		var data setTokenReq
		err = json.Unmarshal(ctx.PostBody(), &data)
		if err != nil {
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

		err = db.SetBotToken(bot_id, token)
		if err != nil {
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
		err = json.Unmarshal(ctx.PostBody(), &data)
		if err != nil {
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

		err = db.SetBotToken(bot_id, &token)
		if err != nil {
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

func startBot(db *pgsql.Db, bots *map[string]*bot.TBot, server *telego.MultiBotWebhookServer, conf *config.BotConfig) fasthttp.RequestHandler {
	// TODO: Check token exist (from db)
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

func healthHandler(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("OK")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func AddHandlers(r *fastRouter.Router, db *pgsql.Db, bots *map[string]*bot.TBot, server *telego.MultiBotWebhookServer, conf *config.BotConfig) {
	r.GET("/api/bot/health", healthHandler)

	r.POST("/api/bot/new", newBot(db))
	r.POST("/api/bot/setToken", setToken(db))

	// Mb change to DELETE http methon
	r.POST("/api/bot/deleteToken", deleteToken(db))

	r.POST("/api/bot/start", startBot(db, bots, server, conf))
}
