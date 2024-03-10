package handlers

import (
	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandler) SetBotToken(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("botId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	data := new(model.SetBotTokenReq)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// check token
	if err := data.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	ok, err := h.bs.TokenHealthCheck(*data.Token)
	if err != nil {
		h.log.Errorw("failed check token health", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrInvalidToken)
	}

	// check bot runnig
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if botStatus == model.StatusBotRunning {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrBotNeedsStopped)
	}

	// check token exists
	existToken, err := h.db.CheckBotTokenExist(data.Token)
	if err != nil {
		h.log.Errorw("failed check bot token exist", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if existToken {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrTokenAlreadyExists)
	}

	if err = h.db.SetBotToken(userId, botId, data.Token); err != nil {
		h.log.Errorw("failed set bot token", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) DeleteBotToken(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("botId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// check bot runnig
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if botStatus == model.StatusBotRunning {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrBotNeedsStopped)
	}

	token := ""

	if err = h.db.SetBotToken(userId, botId, &token); err != nil {
		h.log.Errorw("failed set bot token (delete token)", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) GetBotToken(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("botId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	token, err := h.db.GetBotToken(userId, botId)
	if err != nil {
		h.log.Errorw("failed get bot token", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"token": *token,
	})
}
