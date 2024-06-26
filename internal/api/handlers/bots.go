package handlers

import (
	"github.com/botscubes/bot-components/components"
	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/model"
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
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	data := new(model.NewBotReq)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	title := data.Title
	token := ""

	if err := data.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	m := &model.Bot{
		UserId: userId,
		Token:  &token,
		Title:  title,
		Status: model.StatusBotStopped,
	}

	mc := &model.Component{

		ComponentData: components.ComponentData{
			ComponentTypeData: components.ComponentTypeData{
				Type: components.TypeStart,
			},
		},
		Position: &model.Point{
			X: float64(config.StartComponentPosX), Y: float64(config.StartComponentPosY),
			Valid: true,
		},
	}

	botId, compId, err := h.db.CreateBot(m, mc)
	if err != nil {
		h.log.Errorw("failed create bot", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	mc.Id = compId

	dataRes := &newBotRes{
		BotId:     botId,
		Component: mc,
	}

	return ctx.Status(fiber.StatusCreated).JSON(dataRes)
}
func (h *ApiHandler) DeleteBot(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("BotId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// check bot status is running
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if botStatus == model.StatusBotRunning {
		token, err := h.db.GetBotToken(userId, botId)
		if err != nil {
			h.log.Errorw("failed get bot token", "error", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if token == nil || *token == "" {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrTokenNotFound)
		}

		if err := h.bs.StopBot(*token); err != nil {
			if err.Error() == bot.ErrTgAuth401.Error() {
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrInvalidToken)
			}

			h.log.Errorw("failed stop bot (delete webhook)", "error", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

	}
	if err := h.db.DeleteBot(userId, botId); err != nil {
		h.log.Errorw("failed delete bot", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) StartBot(ctx *fiber.Ctx) error {

	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("BotId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	// check bot already running
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if botStatus == model.StatusBotRunning {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrBotAlreadyRunning)
	}

	token, err := h.db.GetBotToken(userId, botId)
	if err != nil {
		h.log.Errorw("failed get bot token", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if token == nil || *token == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrTokenNotFound)
	}

	// token health check
	ok, err = h.bs.TokenHealthCheck(*token)
	if err != nil {
		h.log.Errorw("failed check token health", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrInvalidToken)
	}

	// starting worker
	if err := h.mb.StartBot(botId, *token); err != nil {
		h.log.Errorw("failed broker: start bot", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// set webhook
	if err = h.bs.StartBot(botId, *token); err != nil {
		if err.Error() == bot.ErrTgAuth401.Error() {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrInvalidToken)
		}

		h.log.Errorw("failed start bot (set webhook)", "error", err)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrStartBot)
	}

	// upd status in db
	if err = h.db.SetBotStatus(botId, userId, model.StatusBotRunning); err != nil {
		h.log.Errorw("failed update bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) StopBot(ctx *fiber.Ctx) error {

	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("BotId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	// check bot already stopped
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if botStatus == model.StatusBotStopped {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrBotAlreadyStopped)
	}

	token, err := h.db.GetBotToken(userId, botId)
	if err != nil {
		h.log.Errorw("failed get bot token", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if token == nil || *token == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrTokenNotFound)
	}

	// delete webhook
	if err := h.bs.StopBot(*token); err != nil {
		if err.Error() == bot.ErrTgAuth401.Error() {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrInvalidToken)
		}

		h.log.Errorw("failed stop bot (delete webhook)", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// update status in db
	if err = h.db.SetBotStatus(botId, userId, model.StatusBotStopped); err != nil {
		h.log.Errorw("failed update bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// stop worker in the background
	go func() {
		if err := h.mb.StopBot(botId); err != nil {
			h.log.Errorw("failed broker: stop bot", "error", err)
			return
		}
	}()

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ApiHandler) GetBots(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	bots, err := h.db.UserBots(userId)
	if err != nil {
		h.log.Errorw("failed get list user bots", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).JSON(bots)
}

func (h *ApiHandler) GetBotStatus(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(int64)
	if !ok {
		h.log.Errorw("UserId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botId, ok := ctx.Locals("botId").(int64)
	if !ok {
		h.log.Errorw("BotId to int64 convert", "error", ErrUserIDConvertation)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	botStatus, err := h.db.GetBotStatus(botId, userId)
	if err != nil {
		h.log.Errorw("failed get bot status", "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).JSON(botStatus)
}
