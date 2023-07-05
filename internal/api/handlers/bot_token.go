package handlers

import (
	"strconv"

	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
)

type setBotTokenReq struct {
	Token *string `json:"token"`
}

func SetBotToken(db *pgsql.Db, log *zap.SugaredLogger) fiber.Handler {
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

		data := new(setBotTokenReq)

		if err := ctx.BodyParser(data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
		}

		// check token
		// TODO: check by tg api
		token := data.Token
		if token == nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.MissingParam("token")))
		}

		if !validateToken(*token) {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrIncorrectTokenFormat))
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

		// check bot already runnig
		botStatus, err := db.GetBotStatus(botId, userId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if botStatus == model.StatusBotRunning {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNeedsStopped))
		}

		// check token exists
		existToken, err := db.CheckBotTokenExist(token)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if existToken {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenAlreadyExists))
		}

		if err = db.SetBotToken(userId, botId, token); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
	}
}

func DeleteBotToken(db *pgsql.Db, log *zap.SugaredLogger) fiber.Handler {
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

		// check bot already runnig
		botStatus, err := db.GetBotStatus(botId, userId)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if botStatus == model.StatusBotRunning {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNeedsStopped))
		}

		token := ""

		if err = db.SetBotToken(userId, botId, &token); err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
	}
}
