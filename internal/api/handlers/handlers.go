package handlers

import (
	"errors"

	e "github.com/botscubes/bot-service/internal/api/errors"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrUserIDConvertation = errors.New("userId convertation to int64")
)

func Health(ctx *fiber.Ctx) error {
	return ctx.SendString("Ok")
}

func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.Status(404).JSON(resp.New(false, nil, e.ErrNotFound))
}
