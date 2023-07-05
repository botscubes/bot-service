package handlers

import (
	"strconv"
	"unicode/utf8"

	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
)

type newBotReq struct {
	Title *string `json:"title"`
}

type newBotRes struct {
	BotId     int64            `json:"botId"`
	Component *model.Component `json:"component"`
}

func NewBot(db *pgsql.Db, log *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		data := new(newBotReq)

		if err := ctx.BodyParser(data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
		}

		title := data.Title
		token := ""

		if title == nil || *title == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.MissingParam("title")))
		}

		if utf8.RuneCountInString(*title) > config.MaxTitleLen {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidTitleLength))
		}

		m := &model.Bot{
			UserId: userId,
			Token:  &token,
			Title:  title,
			Status: model.StatusBotStopped,
		}

		botId, err := db.AddBot(m)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if err := db.CreateBotSchema(botId); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
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
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		mc.Id = compId
		mc.Commands = new(model.Commands)

		dataRes := &newBotRes{
			BotId:     botId,
			Component: mc,
		}

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, dataRes, nil))
	}
}

func StartBot(
	db *pgsql.Db,
	bs *bot.BotService,
	log *zap.SugaredLogger,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if !existBot {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNotFound))
		}

		// check bot already running
		botStatus, err := db.GetBotStatus(botId, userId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if botStatus == model.StatusBotRunning {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotAlreadyRunning))
		}

		token, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if token == nil || *token == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenNotFound))
		}

		// WARN: need sync bs & db bot status (on error)
		if err = bs.StartBot(botId, *token); err != nil {
			if err.Error() == ErrTgAuth401.Error() {
				return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
			}

			log.Error(err)
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrStartBot))
		}

		if err = db.SetBotStatus(botId, userId, model.StatusBotRunning); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
	}
}

func StopBot(db *pgsql.Db, bs *bot.BotService, log *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if !existBot {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNotFound))
		}

		// check bot already stopped
		botStatus, err := db.GetBotStatus(botId, userId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if botStatus == model.StatusBotStopped {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotAlreadyStopped))
		}

		token, err := db.GetBotToken(userId, botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if token == nil || *token == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenNotFound))
		}

		// WARN: need sync bs & db bot status (on error)
		// TODO: error stopping if token invalid
		if err := bs.StopBot(*token); err != nil {
			if err.Error() == ErrTgAuth401.Error() {
				return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
			}

			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrStopBot))
		}

		if err = db.SetBotStatus(botId, userId, model.StatusBotStopped); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
	}
}

func GetBots(db *pgsql.Db, log *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		bots, err := db.UserBots(userId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, bots, nil))
	}
}

func WipeBot(db *pgsql.Db, r *rdb.Rdb, bs *bot.BotService, log *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if !existBot {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNotFound))
		}

		// check bot status is running
		botStatus, err := db.GetBotStatus(botId, userId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if botStatus == model.StatusBotRunning {
			token, err := db.GetBotToken(userId, botId)
			if err != nil {
				log.Error(err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
			}

			if token == nil || *token == "" {
				return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenNotFound))
			}

			if err := bs.StopBot(*token); err != nil {
				if err.Error() == ErrTgAuth401.Error() {
					return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
				}

				log.Error(err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrStopBot))
			}

			if err = db.SetBotStatus(botId, userId, model.StatusBotStopped); err != nil {
				log.Error(err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
			}
		}

		// remove components
		err = db.DelAllComponents(botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		// remove commands
		err = db.DelAllCommands(botId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		// remove next step from main component
		if err = db.DelNextStepComponent(botId, config.MainComponentId); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		// remove token
		token := ""
		if err = db.SetBotToken(userId, botId, &token); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		// Invalidate bot cache
		r.DelBotData(botId)

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
	}
}
