package handlers

import (
	"strconv"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
)

type newBotRes struct {
	BotId     int64            `json:"botId"`
	Component *model.Component `json:"component"`
}

func (h *ApiHandler) NewBot(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	data := new(model.NewBotReq)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBadRequest))
	}

	title := data.Title
	token := ""

	if err := data.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, err))
	}

	m := &model.Bot{
		UserId: userId,
		Token:  &token,
		Title:  title,
		Status: model.StatusBotStopped,
	}

	componentType := "start"

	mc := &model.Component{
		Data: &model.ComponentData{
			Type:    &componentType,
			Content: &[]*model.Content{},
		},
		Keyboard: &model.Keyboard{
			Buttons: [][]*int64{},
		},
		NextStepId: nil,
		IsMain:     true,
		Position: &model.Point{
			X: float64(config.StartComponentPosX), Y: float64(config.StartComponentPosY),
			Valid: true,
		},
		Status: model.StatusComponentActive,
	}

	botId, compId, err := h.db.CreateBot(m, mc)
	if err != nil {
		h.log.Errorw("failed create bot", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	mc.Id = compId
	mc.Commands = new(model.Commands)

	dataRes := &newBotRes{
		BotId:     botId,
		Component: mc,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, dataRes, nil))
}

func (h *ApiHandler) StartBot(ctx *fiber.Ctx) error {
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

	// check bot already running
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if botStatus == model.StatusBotRunning {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotAlreadyRunning))
	}

	token, err := h.db.GetBotToken(userId, botId)
	if err != nil {
		h.log.Errorw("failed get bot token", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if token == nil || *token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenNotFound))
	}

	// token health check
	ok, err = h.bs.TokenHealthCheck(*token)
	if err != nil {
		h.log.Errorw("failed check token health", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
	}

	// starting worker
	if err := h.mb.StartBot(botId, *token); err != nil {
		h.log.Errorw("failed broker: start bot", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// set webhook
	if err = h.bs.StartBot(botId, *token); err != nil {
		if err.Error() == bot.ErrTgAuth401.Error() {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
		}

		h.log.Errorw("failed start bot (set webhook)", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrStartBot))
	}

	// upd status in db
	if err = h.db.SetBotStatus(botId, userId, model.StatusBotRunning); err != nil {
		h.log.Errorw("failed update bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

func (h *ApiHandler) StopBot(ctx *fiber.Ctx) error {
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

	// check bot already stopped
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if botStatus == model.StatusBotStopped {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrBotAlreadyStopped))
	}

	token, err := h.db.GetBotToken(userId, botId)
	if err != nil {
		h.log.Errorw("failed get bot token", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if token == nil || *token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenNotFound))
	}

	// delete webhook
	if err := h.bs.StopBot(*token); err != nil {
		if err.Error() == bot.ErrTgAuth401.Error() {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
		}

		h.log.Errorw("failed stop bot (delete webhook)", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrStopBot))
	}

	// update status in db
	if err = h.db.SetBotStatus(botId, userId, model.StatusBotStopped); err != nil {
		h.log.Errorw("failed update bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// stop worker in the background
	go func() {
		if err := h.mb.StopBot(botId); err != nil {
			h.log.Errorw("failed broker: stop bot", "error", err)
			return
		}
	}()

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}

func (h *ApiHandler) GetBots(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	bots, err := h.db.UserBots(userId)
	if err != nil {
		h.log.Errorw("failed get list user bots", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, bots, nil))
}

func (h *ApiHandler) WipeBot(ctx *fiber.Ctx) error {
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

	// check bot status is running
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	if botStatus == model.StatusBotRunning {
		token, err := h.db.GetBotToken(userId, botId)
		if err != nil {
			h.log.Errorw("failed get bot token", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if token == nil || *token == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrTokenNotFound))
		}

		if err := h.bs.StopBot(*token); err != nil {
			if err.Error() == bot.ErrTgAuth401.Error() {
				return ctx.Status(fiber.StatusBadRequest).JSON(resp.New(false, nil, e.ErrInvalidToken))
			}

			h.log.Errorw("failed stop bot (delete webhook)", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrStopBot))
		}

		if err = h.db.SetBotStatus(botId, userId, model.StatusBotStopped); err != nil {
			h.log.Errorw("failed set bot status", "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}
	}

	// TODO: create transaction
	// remove components
	err = h.db.DelAllComponents(botId)
	if err != nil {
		h.log.Errorw("failed delete all components of bot", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// remove commands
	err = h.db.DelAllCommands(botId)
	if err != nil {
		h.log.Errorw("failed delete all commands of bot", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// remove next step from main component
	if err = h.db.DelNextStepComponent(botId, config.MainComponentId); err != nil {
		h.log.Errorw("failed next step from main component", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// remove token
	token := ""
	if err = h.db.SetBotToken(userId, botId, &token); err != nil {
		h.log.Errorw("failed set bot token (delete token)", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
	}

	// Invalidate bot cache
	h.r.DelBotData(botId)

	return ctx.Status(fiber.StatusOK).JSON(resp.New(true, nil, nil))
}
