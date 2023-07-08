package handlers

import (
	"strconv"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
)

type addCommandReq struct {
	Type *string `json:"type"`
	Data *string `json:"data"`
}

type addCommandRes struct {
	Id int64 `json:"id"`
}

func (h *ApiHandler) AddCommand(ctx *fiber.Ctx) error {
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

	reqData := new(addCommandReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if err := model.ValidateCommand(reqData.Type, reqData.Data); err != nil {
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

	// check bot component exists
	existComp, err := h.db.CheckComponentExist(botId, compId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrComponentNotFound))
	}

	m := &model.Command{
		Type:        reqData.Type,
		Data:        reqData.Data,
		ComponentId: &compId,
		NextStepId:  nil,
		Status:      model.StatusCommandActive,
	}

	commandId, err := h.db.AddCommand(botId, m)
	if err != nil {
		h.log.Errorw("failed add command", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	dataRes := &addCommandRes{
		Id: commandId,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, dataRes, nil))
}

func (h *ApiHandler) DelCommand(ctx *fiber.Ctx) error {
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

	commandId, err := strconv.ParseInt(ctx.Params("commandId"), 10, 64)
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

	// Check command exists
	existCommand, err := h.db.CheckCommandExist(botId, compId, commandId)
	if err != nil {
		h.log.Errorw("failed check command exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existCommand {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrCommandNotFound))
	}

	if err = h.db.DelCommand(botId, commandId); err != nil {
		h.log.Errorw("failed delete command", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

type setNextStepCommandReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func (h *ApiHandler) SetNextStepCommand(ctx *fiber.Ctx) error {
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

	commandId, err := strconv.ParseInt(ctx.Params("commandId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	reqData := new(setNextStepCommandReq)
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

	// Check next component exists
	existNextComp, err := h.db.CheckComponentExist(botId, *nextComponentId)
	if err != nil {
		h.log.Errorw("failed check next component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existNextComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrNextComponentNotFound))
	}

	// Check command exists
	existCommand, err := h.db.CheckCommandExist(botId, compId, commandId)
	if err != nil {
		h.log.Errorw("failed check command exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existCommand {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrCommandNotFound))
	}

	if err = h.db.SetNextStepCommand(botId, commandId, *nextComponentId); err != nil {
		h.log.Errorw("failed set next step command", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

func (h *ApiHandler) DelNextStepCommand(ctx *fiber.Ctx) error {
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

	commandId, err := strconv.ParseInt(ctx.Params("commandId"), 10, 64)
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

	// Check command exists
	existCommand, err := h.db.CheckCommandExist(botId, compId, commandId)
	if err != nil {
		h.log.Errorw("failed check command exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existCommand {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrCommandNotFound))
	}

	if err = h.db.DelNextStepCommand(botId, commandId); err != nil {
		h.log.Errorw("failed delete next step from command", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

type updCommandReq struct {
	Type *string `json:"type"`
	Data *string `json:"data"`
}

func (h *ApiHandler) UpdCommand(ctx *fiber.Ctx) error {
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

	commandId, err := strconv.ParseInt(ctx.Params("commandId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	reqData := new(updCommandReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if err := model.ValidateCommand(reqData.Type, reqData.Data); err != nil {
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

	// check bot component exists
	existComp, err := h.db.CheckComponentExist(botId, compId)
	if err != nil {
		h.log.Errorw("failed check component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrComponentNotFound))
	}

	// Check command exists
	existCommand, err := h.db.CheckCommandExist(botId, compId, commandId)
	if err != nil {
		h.log.Errorw("failed check command exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existCommand {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrCommandNotFound))
	}

	err = h.db.UpdCommand(botId, commandId, reqData.Type, reqData.Data)
	if err != nil {
		h.log.Errorw("failed update command", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}
