package middlewares

import (
	"strconv"

	"github.com/botscubes/bot-service/internal/api/handlers"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
)

func GetGroupMiddleware(db *pgsql.Db, log *zap.SugaredLogger,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		botId, ok := ctx.Locals("botId").(int64)
		if !ok {
			log.Errorw("botId to int64 convert", "error", handlers.ErrUserIDConvertation)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		groupId, err := strconv.ParseInt(ctx.Params("groupId"), 10, 64)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		existGroup, err := db.CheckGroupExist(botId, groupId)
		if err != nil {
			log.Errorw("failed check group exist", "error", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if !existGroup {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrGroupNotFound)
		}

		ctx.Locals("groupId", botId)

		return ctx.Next()
	}
}
