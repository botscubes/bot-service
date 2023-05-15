package handlers

import (
	"strconv"
	"unicode/utf8"

	"github.com/goccy/go-json"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	fh "github.com/valyala/fasthttp"
)

type newBotReq struct {
	Title *string `json:"title"`
}

type newBotRes struct {
	BotId     int64            `json:"botId"`
	Component *model.Component `json:"component"`
}

func NewBot(db *pgsql.Db, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		var data newBotReq

		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		title := data.Title
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		token := ""

		if title == nil || *title == "" {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.MissingParam("title")))
			return
		}

		if utf8.RuneCountInString(*title) > config.MaxTitleLen {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidTitleLength))
			return
		}

		m := &model.Bot{
			UserId: userId,
			Token:  &token,
			Title:  title,
			Status: model.StatusBotActive,
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

		mc := &model.Component{
			Data: &model.Data{
				Type:    &dataType,
				Content: &[]*model.Content{},
			},
			Keyboard: &model.Keyboard{
				Buttons: [][]*int64{},
			},
			NextStepId: nil,
			IsMain:     true,
			Position: &model.Point{
				X: float64(config.StartComponentPosX), Y: float64(config.StartComponentPosY),
				Valid: true,
			},
			Status: model.StatusComponentActive,
		}

		compId, err := db.AddComponent(botId, mc)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		mc.Id = compId
		mc.Commands = new(model.Commands)

		dataRes := &newBotRes{
			BotId:     botId,
			Component: mc,
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, dataRes, nil))
	}
}

func StartBot(
	db *pgsql.Db,
	bs *bot.BotService,
	r *rdb.Rdb,
	log *zap.SugaredLogger,
) reqHandler {
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

		token, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if token == nil || *token == "" {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrTokenNotFound))
			return
		}

		if ok := bs.CheckBotExist(botId); !ok {
			// TODO: Own token health check to get a specific error
			if err = bs.NewBot(token, botId, log, r, db); err != nil {
				doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidToken))
				return
			}
		}

		// check bot already runnig
		isRunning, err := bs.BotIsRunnig(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if isRunning {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotAlreadyRunning))
			return
		}

		if err = bs.StartBot(botId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrStartBot))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func StopBot(db *pgsql.Db, bs *bot.BotService, log *zap.SugaredLogger) reqHandler {
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

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotAlreadyRunning))
			return
		}

		token, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if token == nil || *token == "" {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrTokenNotFound))
			return
		}

		if ok := bs.CheckBotExist(botId); !ok {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotAlreadyStopped))
			return
		}

		// check bot already stopped
		isRunning, err := bs.BotIsRunnig(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !isRunning {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotAlreadyStopped))
			return
		}

		if err := bs.StopBot(botId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrStopBot))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func GetBots(db *pgsql.Db, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		bots, err := db.UserBots(userId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, bots, nil))
	}
}

func WipeBot(db *pgsql.Db, bs *bot.BotService, log *zap.SugaredLogger) reqHandler {
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
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotAlreadyRunning))
			return
		}

		// stop bot worker
		if ok := bs.CheckBotExist(botId); ok {
			// check bot already stopped
			isRunning, err := bs.BotIsRunnig(botId)
			if err != nil {
				log.Error(err)
				doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
				return
			}

			if isRunning {
				if err := bs.StopBot(botId); err != nil {
					log.Error(err)
					doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrStopBot))
					return
				}
			}
		}

		// remove components
		err = db.DelAllComponents(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// remove commands
		err = db.DelAllCommands(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// remove token
		token := ""
		if err = db.SetBotToken(userId, botId, &token); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}
