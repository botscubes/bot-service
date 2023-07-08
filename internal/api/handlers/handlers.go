package handlers

import (
	"errors"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	ErrUserIDConvertation = errors.New("userId convertation to int64")
)

type ApiHandler struct {
	db  *pgsql.Db
	log *zap.SugaredLogger
	bs  *bot.BotService
	nc  *nats.Conn
	r   *rdb.Rdb
}

func NewApiHandler(
	db *pgsql.Db,
	log *zap.SugaredLogger,
	bs *bot.BotService,
	nc *nats.Conn,
	r *rdb.Rdb,
) *ApiHandler {
	return &ApiHandler{
		db:  db,
		log: log,
		bs:  bs,
		nc:  nc,
		r:   r,
	}
}

func Health(ctx *fiber.Ctx) error {
	return ctx.SendString("Ok")
}

func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(resp.New(false, nil, e.ErrNotFound))
}
