package middlewares

import (
	"context"
	"strings"

	"github.com/botscubes/user-service/pkg/jwt"
	"github.com/botscubes/user-service/pkg/token_storage"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Auth(st *token_storage.TokenStorage, jwtKey *string, log *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		const prefix = "Bearer "

		auth := ctx.Get("Authorization")

		if !strings.HasPrefix(auth, prefix) {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		token := strings.TrimPrefix(auth, prefix)
		exists, err := (*st).CheckToken(context.Background(), token)
		if err != nil {
			log.Errorw("failed check JWT exists (auth)", "errror", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if !exists {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		// WARN: fix error !!!
		userId, err := jwt.GetIdFromToken(token, *jwtKey)
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		ctx.Locals("userId", int64(userId))

		return ctx.Next()
	}
}
