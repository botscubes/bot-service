package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandler) GetBotComponents(ctx *fiber.Ctx) error {
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("BotId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	groupId, ok := ctx.Locals("groupId").(int64)
	if !ok {
		h.log.Errorw("GroupId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
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
