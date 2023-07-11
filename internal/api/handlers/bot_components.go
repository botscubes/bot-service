package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
)

type AddComponentRes struct {
	Id int64 `json:"id"`
}

func (h *ApiHandler) AddComponent(ctx *fiber.Ctx) error {
	var err error
	botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	reqData := new(model.AddComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if errValidate := reqData.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, errValidate))
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

	tx, err := h.db.BeginTx(ctx.Context())
	if err != nil {
		h.log.Errorw("failed begin db transaction", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx.Context())
		} else {
			_ = tx.Commit(ctx.Context())
		}
	}()

	compId, err := h.db.AddComponentTx(ctx.Context(), tx, botId, component)
	if err != nil {
		h.log.Errorw("failed add component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if reqData.Commands != nil {
		for _, v := range *reqData.Commands {
			mc := &model.Command{
				Type:        v.Type,
				Data:        v.Data,
				ComponentId: &compId,
				NextStepId:  nil,
				Status:      model.StatusCommandActive,
			}

			_, err = h.db.AddCommandTx(ctx.Context(), tx, botId, mc)
			if err != nil {
				h.log.Errorw("failed add command", "error", err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
			}
		}
	}

	dataRes := &AddComponentRes{
		Id: compId,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, dataRes, nil))
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

	reqData := new(model.SetNextStepComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if err := reqData.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, err))
	}

	if *reqData.NextStepId == compId {
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
	existNextComp, err := h.db.CheckComponentExist(botId, *reqData.NextStepId)
	if err != nil {
		h.log.Errorw("failed check next component exist", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !existNextComp {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrNextComponentNotFound))
	}

	if err = h.db.SetNextStepComponent(botId, compId, *reqData.NextStepId); err != nil {
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
	var err error
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

	tx, err := h.db.BeginTx(ctx.Context())
	if err != nil {
		h.log.Errorw("failed begin db transaction", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx.Context())
		} else {
			_ = tx.Commit(ctx.Context())
		}
	}()

	// delete component
	if err = h.db.DelComponentTx(ctx.Context(), tx, botId, compId); err != nil {
		h.log.Errorw("failed delete component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// delete component commands
	if err = h.db.DelCommandsByCompIdTx(ctx.Context(), tx, botId, compId); err != nil {
		h.log.Errorw("failed delete commands by component id", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// delete component next steps, that reference these component
	if err = h.db.DelNextStepComponentByNsTx(ctx.Context(), tx, botId, compId); err != nil {
		h.log.Errorw("failed delete component next steps, that reference these component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate component cache
	if err = h.r.DelComponent(botId, compId); err != nil {
		h.log.Errorw("failed delete component from cache", "error", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
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

	reqData := new(model.DelSetComponentsReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if err := reqData.Validate(); err != nil {
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

	reqData := new(model.UpdComponentReq)
	if err := ctx.BodyParser(reqData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	if err := reqData.Validate(); err != nil {
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

	delCache := false

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
		delCache = true
	}

	if delCache {
		// Invalidate component cache
		if err = h.r.DelComponent(botId, compId); err != nil {
			h.log.Errorw("failed delete component from cache", "error", err)
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}
