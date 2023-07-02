package middlewares

import (
	"context"
	"strings"

	e "github.com/botscubes/bot-service/internal/api/errors"
	resp "github.com/botscubes/bot-service/pkg/api_response"

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
			return ctx.Status(fiber.StatusUnauthorized).JSON(resp.New(false, nil, e.ErrUnauthorized))
		}

		token := strings.TrimPrefix(auth, prefix)
		exists, err := (*st).CheckToken(context.Background(), token)
		if err != nil {
			log.Error(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(resp.New(false, nil, e.ErrInternalServer))
		}

		if !exists {
			return ctx.Status(fiber.StatusUnauthorized).JSON(resp.New(false, nil, e.ErrUnauthorized))
		}

		// WARN: fix error !!!
		userId, err := jwt.GetIdFromToken(token, *jwtKey)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(resp.New(false, nil, e.ErrUnauthorized))
		}

		ctx.Locals("userId", int64(userId))

		return ctx.Next()
	}
}
