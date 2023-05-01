package app

import (
	h "github.com/botscubes/bot-service/internal/api/handlers"
)

func (app *App) addHandlers() {
	app.Router.GET("/api/bots/health", h.Auth(h.Health, &app.SessionStorage, &app.Conf.JWTKey))

	// Create new bot
	app.Router.POST("/api/bots/new", h.Auth(h.NewBot(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Set bot token
	app.Router.POST("/api/bots/{botId}/token", h.Auth(h.SetBotToken(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Delete bot token
	app.Router.DELETE("/api/bots/{botId}/token", h.Auth(h.DeleteBotToken(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Start bot
	app.Router.PATCH("/api/bots/{botId}/start", h.Auth(h.StartBot(app.Db, &app.Bots, app.Server, &app.Conf.Bot), &app.SessionStorage, &app.Conf.JWTKey))

	// Stop bot
	app.Router.PATCH("/api/bots/{botId}/stop", h.Auth(h.StopBot(app.Db, &app.Bots), &app.SessionStorage, &app.Conf.JWTKey))

	// Bot components

	// Adds a component to the bot structure
	app.Router.POST("/api/bots/{botId}/components/add", h.Auth(h.AddBotComponent(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Set next step from the component
	app.Router.POST("/api/bots/{botId}/components/{compId}/next", h.Auth(h.SetNextForComponent(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Set next step from the component command
	app.Router.POST("/api/bots/{botId}/components/{compId}/commands/{commandId}/next", h.Auth(h.SetNextForCommand(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Get bot components
	app.Router.GET("/api/bots/{botId}/components", h.Auth(h.GetBotComponents(app.Db), &app.SessionStorage, &app.Conf.JWTKey))
}
