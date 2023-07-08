package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
)

type addComponentReq struct {
	Data     *model.Data     `json:"data"`
	Commands *model.Commands `json:"commands"`
	Position *model.Point    `json:"position"`
}

type addComponentRes struct {
	Id int64 `json:"id"`
}

func (h *ApiHandler) AddComponent(ctx *fiber.Ctx) error {
	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	reqData := new(addComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	// TODO: check fields limits:
	// eg. data.commands._.data max size, check commands max count
	if err := model.ValidateComponent(reqData.Data, reqData.Commands, reqData.Position); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, err))
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

	component := &model.Component{
		Data: reqData.Data,
		Keyboard: &model.Keyboard{
			Buttons: [][]*int64{},
		},
		NextStepId: nil,
		IsMain:     false,
		Position:   reqData.Position,
		Status:     model.StatusComponentActive,
	}

	compId, err := h.db.AddComponent(botId, component)
	if err != nil {
		h.log.Errorw("failed add component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	for _, v := range *reqData.Commands {
		mc := &model.Command{
			Type:        v.Type,
			Data:        v.Data,
			ComponentId: &compId,
			NextStepId:  nil,
			Status:      model.StatusCommandActive,
		}

		_, err := h.db.AddCommand(botId, mc)
		if err != nil {
			h.log.Errorw("failed add command", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}
	}

	dataRes := &addComponentRes{
		Id: compId,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, dataRes, nil))
}

type setNextStepComponentReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func (h *ApiHandler) SetNextStepComponent(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	compId, err := strconv.ParseInt(ctx.Params("compId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	reqData := new(setNextStepComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if reqData.NextStepId == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.MissingParam("nextStepId")))
	}

	nextComponentId := reqData.NextStepId

	if *nextComponentId == compId {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.InvalidParam("nextStepId")))
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

	// check bot component exists
	existInitialComp, err := h.db.CheckComponentExist(botId, compId)
	if err != nil {
		h.log.Errorw("failed check initial component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existInitialComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrComponentNotFound))
	}

	// check bot next component exists
	existNextComp, err := h.db.CheckComponentExist(botId, *nextComponentId)
	if err != nil {
		h.log.Errorw("failed check next component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existNextComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrNextComponentNotFound))
	}

	if err = h.db.SetNextStepComponent(botId, compId, *nextComponentId); err != nil {
		h.log.Errorw("failed set next step for component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

func (h *ApiHandler) GetBotComponents(ctx *fiber.Ctx) error {
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

	components, err := h.db.ComponentsForEd(botId)
	if err != nil {
		h.log.Errorw("failed get bot components for editor", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, components, nil))
}

func (h *ApiHandler) DelNextStepComponent(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	compId, err := strconv.ParseInt(ctx.Params("compId"), 10, 64)
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

	// check bot component exists
	existComp, err := h.db.CheckComponentExist(botId, compId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrComponentNotFound))
	}

	if err = h.db.DelNextStepComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete next step from component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

func (h *ApiHandler) DelComponent(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	compId, err := strconv.ParseInt(ctx.Params("compId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	// check component is main
	if compId == config.MainComponentId {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrMainComponent))
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

	// check bot component exists
	existComp, err := h.db.CheckComponentExist(botId, compId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrComponentNotFound))
	}

	if err = h.db.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// delete component commands
	if err = h.db.DelCommandsByCompId(botId, compId); err != nil {
		h.log.Errorw("failed delete commands by component id", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// delete component next steps, that reference these component
	if err = h.db.DelNextStepComponentByNS(botId, compId); err != nil {
		h.log.Errorw("failed delete component next steps, that reference these component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

type delSetComponentsReq struct {
	Data *[]int64 `json:"data"`
}

func (h *ApiHandler) DelSetOfComponents(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	reqData := new(delSetComponentsReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if reqData.Data == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if len(*reqData.Data) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidParam))
	}

	// check exist main component
	for _, v := range *reqData.Data {
		if v == config.MainComponentId {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrMainComponent))
		}
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

	// delete components
	if err = h.db.DelSetOfComponents(botId, reqData.Data); err != nil {
		h.log.Errorw("failed delete set of components", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	for _, v := range *reqData.Data {
		if err = h.r.DelComponent(botId, v); err != nil {
			h.log.Errorw("failed delete component from cache", "error", err)
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

type updComponentReq struct {
	Data     *model.Data  `json:"data"`
	Position *model.Point `json:"position"`
}

func (h *ApiHandler) UpdComponent(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	compId, err := strconv.ParseInt(ctx.Params("compId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	reqData := new(updComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if reqData.Data == nil && reqData.Position == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if reqData.Data != nil {
		if err := reqData.Data.Validate(); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, err))
		}
	}

	if reqData.Position != nil {
		if err := reqData.Position.Validate(); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, err))
		}
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

	// check bot component exists
	existComp, err := h.db.CheckComponentExist(botId, compId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrComponentNotFound))
	}

	if reqData.Position != nil {
		err = h.db.UpdComponentPosition(botId, compId, reqData.Position)
		if err != nil {
			h.log.Errorw("failed update component position", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}
	}

	if reqData.Data != nil {
		err = h.db.UpdComponentData(botId, compId, reqData.Data)
		if err != nil {
			h.log.Errorw("failed update component data", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}
