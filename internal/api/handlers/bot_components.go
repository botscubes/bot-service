package handlers

import (
	"github.com/botscubes/bot-components/components"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/gofiber/fiber/v2"

	e "github.com/botscubes/bot-service/internal/api/errors"
)

type AddComponentRes struct {
	Id int64 `json:"id"`
}

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

func (h *ApiHandler) AddComponent(ctx *fiber.Ctx) error {
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

	reqData := new(model.AddComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if errValidate := reqData.Validate(); errValidate != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errValidate)
	}

	compId, err := h.db.AddComponent(botId, groupId, &model.Component{
		Position: reqData.Position,
		ComponentData: components.ComponentData{
			ComponentTypeData: components.ComponentTypeData{
				Type: reqData.Type,
			},
			Path: reqData.Type,
		},
	})
	if err != nil {
		h.log.Errorw("failed add component", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	dataRes := &AddComponentRes{
		Id: compId,
	}
	return ctx.Status(fiber.StatusCreated).JSON(dataRes)
}

func (h *ApiHandler) DeleteComponent(ctx *fiber.Ctx) error {
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

	componentId, ok := ctx.Locals("componentId").(int64)
	if !ok {
		h.log.Errorw("GroupId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if componentId == 1 {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrDeleteStartComponent)
	}

	err := h.db.DeleteComponent(botId, groupId, componentId)
	if err != nil {
		h.log.Errorw("failed delete component", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) SetComponentPosition(ctx *fiber.Ctx) error {
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

	componentId, ok := ctx.Locals("componentId").(int64)
	if !ok {
		h.log.Errorw("GroupId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	position := new(model.Point)
	if err := ctx.BodyParser(position); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if errValidate := position.Validate(); errValidate != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errValidate)
	}
	err := h.db.SetComponentPosition(botId, groupId, componentId, position)
	if err != nil {
		h.log.Errorw("failed set component position", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) UpdateComponentData(ctx *fiber.Ctx) error {
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

	componentId, ok := ctx.Locals("componentId").(int64)
	if !ok {
		h.log.Errorw("GroupId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	componentType, err := h.db.GetComponentType(botId, groupId, componentId)
	if err != nil {
		h.log.Errorw("failed get component type", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	data := new(map[string]any)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if errValidate := model.ValidateSpecificComponentData(componentType, *data); errValidate != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errValidate)
	}

	if err = h.db.UpdateComponentData(botId, groupId, componentId, *data); err != nil {

		h.log.Errorw("failed set component data", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
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
