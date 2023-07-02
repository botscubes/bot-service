package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrUserIDConvertation = errors.New("userId convertation to int64")
)

func Health(ctx *fiber.Ctx) error {
	return ctx.SendString("Ok")
}
