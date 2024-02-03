package handlers

import (
	"errors"

	"github.com/botscubes/bot-service/internal/bot"
	mb "github.com/botscubes/bot-service/internal/broker"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	ErrUserIDConvertation = errors.New("userId convertation to int64")
)

type ApiHandler struct {
	db  *pgsql.Db
	log *zap.SugaredLogger
	bs  *bot.BotService
	mb  mb.Broker
	r   *rdb.Rdb
}

func NewApiHandler(
	db *pgsql.Db,
	log *zap.SugaredLogger,
	bs *bot.BotService,
	b mb.Broker,
	r *rdb.Rdb,
) *ApiHandler {
	return &ApiHandler{
		db:  db,
		log: log,
		bs:  bs,
		mb:  b,
		r:   r,
	}
}

func Health(ctx *fiber.Ctx) error {
	return ctx.SendString("Ok")
}

func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotFound)
}
