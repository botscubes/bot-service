package app

import (
	h "github.com/botscubes/bot-service/internal/api/handlers"
)

// Bot handlers
func (app *App) regBotsHandlers() {
	// Create new bot
	app.Router.POST("/api/bots",
		h.Auth(
			h.NewBot(app.Db, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Set bot token
	app.Router.POST("/api/bots/{botId}/token",
		h.Auth(
			h.SetBotToken(app.Db, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Delete bot token
	app.Router.DELETE("/api/bots/{botId}/token",
		h.Auth(
			h.DeleteBotToken(app.Db, app.Log, app.BotService),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Start bot
	app.Router.PATCH("/api/bots/{botId}/start",
		h.Auth(
			h.StartBot(app.Db, app.BotService, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Stop bot
	app.Router.PATCH("/api/bots/{botId}/stop",
		h.Auth(
			h.StopBot(app.Db, app.BotService, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Get user bots
	app.Router.GET("/api/bots",
		h.Auth(
			h.GetBots(app.Db, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Wipe bot data
	app.Router.PATCH("/api/bots/{botId}/wipe",
		h.Auth(
			h.WipeBot(app.Db, app.BotService, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))
}

func (app *App) regComponentsHandlers() {
	// Adds a component to the bot structure
	app.Router.POST("/api/bots/{botId}/components",
		h.Auth(h.AddComponent(app.Db, app.Log), &app.SessionStorage, &app.Conf.JWTKey, app.Log))

	// Delete bot component
	app.Router.DELETE("/api/bots/{botId}/components/{compId}",
		h.Auth(
			h.DelComponent(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// update component
	app.Router.PATCH("/api/bots/{botId}/components/{compId}",
		h.Auth(
			h.UpdComponent(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Set next step component
	app.Router.POST("/api/bots/{botId}/components/{compId}/next",
		h.Auth(
			h.SetNextStepComponent(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Get bot components
	app.Router.GET("/api/bots/{botId}/components",
		h.Auth(
			h.GetBotComponents(app.Db, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Delete next step component
	app.Router.DELETE("/api/bots/{botId}/components/{compId}/next",
		h.Auth(
			h.DelNextStepComponent(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log))
}

func (app *App) regCommandsHandlers() {
	// Add command
	app.Router.POST("/api/bots/{botId}/components/{compId}/commands",
		h.Auth(
			h.AddCommand(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Delete component command
	app.Router.DELETE("/api/bots/{botId}/components/{compId}/commands/{commandId}",
		h.Auth(
			h.DelCommand(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// update command
	app.Router.PATCH("/api/bots/{botId}/components/{compId}/commands/{commandId}",
		h.Auth(
			h.UpdCommand(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Set next step component command
	app.Router.POST("/api/bots/{botId}/components/{compId}/commands/{commandId}/next",
		h.Auth(
			h.SetNextStepCommand(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))

	// Delete next step command
	app.Router.DELETE("/api/bots/{botId}/components/{compId}/commands/{commandId}/next",
		h.Auth(
			h.DelNextStepCommand(app.Db, app.Redis, app.Log),
			&app.SessionStorage, &app.Conf.JWTKey, app.Log,
		))
}

func (app *App) regiterHandlers() {
	app.Router.PanicHandler = h.PanicHandler(app.Log)

	app.Router.GET("/api/bots/health",
		h.Auth(h.Health, &app.SessionStorage, &app.Conf.JWTKey, app.Log))

	app.regBotsHandlers()
	app.regComponentsHandlers()
	app.regCommandsHandlers()
}
