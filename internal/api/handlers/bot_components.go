package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	e "github.com/botscubes/bot-service/internal/api/errors"
)

func (h *ApiHandler) GetBotComponents(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// check bot exists
	existBot, err := h.db.CheckBotExist(userId, botId)
	if err != nil {
		h.log.Errorw("failed check bot exist", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !existBot {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrBotNotFound)
	}

	groupId, err := strconv.ParseInt(ctx.Params("groupId"), 10, 64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	existGroup, err := h.db.CheckGroupExist(botId, groupId)
	if err != nil {
		h.log.Errorw("failed check group exist", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if !existGroup {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrGroupNotFound)
	}
	components, err := h.db.GetComponents(botId, groupId)
	if err != nil {
		h.log.Errorw("failed get bot components for editor", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).JSON(components)
}

//func (h *ApiHandler) SetPosition(ctx *fiber.Ctx) error {
//	userId, ok := ctx.Locals("userId").(int64)
//	if !ok {
//		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
//	if err != nil {
//		return ctx.SendStatus(fiber.StatusBadRequest)
//	}
//}
