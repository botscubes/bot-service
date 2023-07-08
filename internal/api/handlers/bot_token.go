package handlers

import (
	"strconv"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
)

type setBotTokenReq struct {
	Token *string `json:"token"`
}

func (h *ApiHandler) SetBotToken(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
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
	existBot, err := h.db.CheckBotExist(userId, botId)
	if err != nil {
		h.log.Errorw("failed check bot exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existBot {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNotFound))
	}

	// check bot already runnig
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if botStatus == model.StatusBotRunning {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNeedsStopped))
	}

	// check token exists
	existToken, err := h.db.CheckBotTokenExist(token)
	if err != nil {
		h.log.Errorw("failed check bot token exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if existToken {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenAlreadyExists))
	}

	if err = h.db.SetBotToken(userId, botId, token); err != nil {
		h.log.Errorw("failed set bot token", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

func (h *ApiHandler) DeleteBotToken(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	// check bot exists
	existBot, err := h.db.CheckBotExist(userId, botId)
	if err != nil {
		h.log.Errorw("failed check bot exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existBot {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNotFound))
	}

	// check bot already runnig
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if botStatus == model.StatusBotRunning {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotNeedsStopped))
	}

	token := ""

	if err = h.db.SetBotToken(userId, botId, &token); err != nil {
		h.log.Errorw("failed set bot token (delete token)", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}
