package middlewares

import (
	"strconv"

	"github.com/botscubes/bot-service/internal/api/handlers"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
)

func GetComponentMiddleware(db *pgsql.Db, log *zap.SugaredLogger,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		botId, ok := ctx.Locals("botId").(int64)
		if !ok {
			log.Errorw("botId to int64 convert", "error", handlers.ErrUserIDConvertation)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		groupId, ok := ctx.Locals("groupId").(int64)
		if !ok {
			log.Errorw("groupId to int64 convert", "error", handlers.ErrUserIDConvertation)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		componentId, err := strconv.ParseInt(ctx.Params("componentId"), 10, 64)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		existComponent, err := db.CheckComponentExist(botId, groupId, componentId)
		if err != nil {
			log.Errorw("failed check component exist", "error", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if !existComponent {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(e.ErrComponentNotFound)
		}
		ctx.Locals("componentId", componentId)

		return ctx.Next()
	}
}
