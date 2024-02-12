package handlers

import (
	"github.com/botscubes/bot-service/internal/model"
	"github.com/gofiber/fiber/v2"

	e "github.com/botscubes/bot-service/internal/api/errors"
)

func (h *ApiHandler) AddConnetion(ctx *fiber.Ctx) error {
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

	reqData := new(model.Connection)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	componentExist, err := h.db.CheckComponentExist(botId, groupId, *reqData.SourceComponentId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !componentExist {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrComponentNotFound)
	}
	componentExist, err = h.db.CheckComponentExist(botId, groupId, *reqData.TargetComponentId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !componentExist {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrComponentNotFound)
	}

	componentType, err := h.db.GetComponentType(botId, groupId, *reqData.SourceComponentId)
	if err != nil {
		h.log.Errorw("failed get component type", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if errValidate := reqData.Validate(componentType); errValidate != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errValidate)
	}

	err = h.db.AddConnection(botId, groupId, reqData)
	if err != nil {
		h.log.Errorw("failed add connection", "error", err)

		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *ApiHandler) DeleteConnection(ctx *fiber.Ctx) error {
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

	reqData := new(model.SourceConnectionPoint)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	componentExist, err := h.db.CheckComponentExist(botId, groupId, *reqData.SourceComponentId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !componentExist {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrComponentNotFound)
	}

	componentType, err := h.db.GetComponentType(botId, groupId, *reqData.SourceComponentId)
	if err != nil {
		h.log.Errorw("failed get component type", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if errValidate := reqData.Validate(componentType); errValidate != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errValidate)
	}

	targetComponentId, err := h.db.GetTargetComponentId(botId, groupId, reqData)
	if err != nil {
		h.log.Errorw("failed get target component id", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)

	}
	if targetComponentId == nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrTargetComponentIdIsNull)
	}

	err = h.db.DeleteConnection(botId, groupId, reqData, *targetComponentId)
	if err != nil {
		h.log.Errorw("failed delete connection", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)

}
