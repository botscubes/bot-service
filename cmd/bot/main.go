package main

import (
	"github.com/botscubes/bot-service/internal/api/handlers"
	"github.com/botscubes/bot-service/internal/app"
)

func main() {
	app := app.New()
	handlers.AddHandlers(app)
	app.Run()
}
