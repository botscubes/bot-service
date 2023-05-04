package handlers

import (
	"strconv"
	"unicode/utf8"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/pgtype"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
	fh "github.com/valyala/fasthttp"
)

type newBotReq struct {
	Title *string `json:"title"`
}

type newBotRes struct {
	BotId     int64      `json:"botId"`
	Component *component `json:"conponent"`
}

func NewBot(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		var data newBotReq

		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			log.Debug("[API: newBot] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		title := data.Title
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: newBot] - get userId convertation to int64 error")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		token := ""
		status := 0

		if title == nil || *title == "" {
			log.Debug("[API: newBot] - title is misssing")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidParams))
			return
		}

		if utf8.RuneCountInString(*title) > config.MaxTitleLen {
			log.Debug("[API: newBot] - title len > ", config.MaxTitleLen)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidTitleLength))
			return
		}

		m := &model.Bot{
			UserId: userId,
			Token:  &token,
			Title:  title,
			Status: status,
		}

		botId, err := db.AddBot(m)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if err := db.CreateBotSchema(botId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		dataType := "start"
		px := 50
		py := 50

		mc := &model.Component{
			Data: &model.Data{
				Type:    &dataType,
				Content: nil,
			},
			Keyboard: &model.Keyboard{
				Buttons: [][]*int64{},
			},
			NextStepId: nil,
			IsStart:    true,
			Position: &pgtype.Point{
				P:      pgtype.Vec2{X: float64(px), Y: float64(py)},
				Status: pgtype.Present,
			},
			Status: pgsql.StatusComponentActive,
		}

		compId, err := db.AddBotComponent(botId, mc)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		mc.Id = compId
		mc.Commands = []*model.Command{}

		dataRes := &newBotRes{
			BotId:     botId,
			Component: componentRes(mc),
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, dataRes, nil))
	}
}

func StartBot(db *pgsql.Db, bots *map[string]*bot.TBot, s *telego.MultiBotWebhookServer, c *config.BotConfig) reqHandler {
	// check bot already started
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: StartBot] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: StartBot] - get userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: startBot] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		token, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if token == nil || *token == "" {
			log.Debug("[API: startBot] - Token not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrTokenNotFound))
			return
		}

		if _, ok := (*bots)[*token]; !ok {
			// TODO: Own token health check to get a specific error
			nbot, err := bot.NewBot(token)
			if err != nil {
				log.Debug("[API: startBot]\n", err)
				doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidToken))
				return
			}

			(*bots)[*token] = new(bot.TBot)
			(*bots)[*token].Bot = nbot
		}

		if err = (*bots)[*token].StartBot(c.WebhookBase, c.ListenAddress, s); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrStartBot))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func StopBot(db *pgsql.Db, bots *map[string]*bot.TBot) reqHandler {
	// TODO: check bot is running
	// check bot already stopped
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: StopBot] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: StopBot] - get userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: stopBot] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		token, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if token == nil || *token == "" {
			log.Debug("[API: stopBot] - Token not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrTokenNotFound))
			return
		}

		if _, ok := (*bots)[*token]; !ok {
			log.Debug("[API: startBot] Bot not found in bots map")
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrBotNotRunning))
			return
		}

		if err := (*bots)[*token].StopBot(false); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrStopBot))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}
