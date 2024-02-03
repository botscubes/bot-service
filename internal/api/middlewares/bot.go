package middlewares

import (
	"strconv"

	"github.com/botscubes/bot-service/internal/api/handlers"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
)

func GetBotMiddleware(db *pgsql.Db, log *zap.SugaredLogger,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int64)
		if !ok {
			log.Errorw("UserId to int64 convert", "error", handlers.ErrUserIDConvertation)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		botId, err := strconv.ParseInt(ctx.Params("botId"), 10, 64)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Errorw("failed check bot exist", "error", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if !existBot {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrBotNotFound)
		}
		ctx.Locals("botId", botId)

		return ctx.Next()
	}
}
