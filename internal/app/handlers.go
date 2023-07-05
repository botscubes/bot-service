package app

import (
	h "github.com/botscubes/bot-service/internal/api/handlers"
	m "github.com/botscubes/bot-service/internal/api/middlewares"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (app *App) regiterHandlers() {
	// Middlewares
	app.Server.Use(recover.New())
	app.Server.Use(m.Auth(&app.SessionStorage, &app.Conf.JWTKey, app.Log))

	app.Server.Get("/api/bots/health", h.Health)

	app.regBotsHandlers()
	app.regComponentsHandlers()
	app.regCommandsHandlers()

	// custom 404 handler
	app.Server.Use(h.NotFoundHandler)
}

// Bot handlers
func (app *App) regBotsHandlers() {
	// Create new bot
	app.Server.Post("/api/bots", h.NewBot(app.Db, app.Log))

	// Set bot token
	app.Server.Post("/api/bots/:botId<int>/token", h.SetBotToken(app.Db, app.Log))

	// Delete bot token
	app.Server.Delete("/api/bots/:botId<int>/token", h.DeleteBotToken(app.Db, app.Log))

	// Start bot
	app.Server.Patch("/api/bots/:botId<int>/start", h.StartBot(app.Db, app.BotService, app.Log))

	// Stop bot
	app.Server.Patch("/api/bots/:botId<int>/stop", h.StopBot(app.Db, app.BotService, app.Log))

	// Get user bots
	app.Server.Get("/api/bots", h.GetBots(app.Db, app.Log))

	// Wipe bot data
	app.Server.Patch("/api/bots/:botId<int>/wipe", h.WipeBot(app.Db, app.Redis, app.BotService, app.Log))
}

func (app *App) regComponentsHandlers() {
	// Adds a component to the bot structure
	app.Server.Post("/api/bots/:botId<int>/components", h.AddComponent(app.Db, app.Log))

	// Delete bot component
	app.Server.Delete("/api/bots/:botId<int>/components/:compId<int>", h.DelComponent(app.Db, app.Redis, app.Log))

	// Delete set of components
	app.Server.Post("/api/bots/:botId<int>/components/del", h.DelSetOfComponents(app.Db, app.Redis, app.Log))

	// update component
	app.Server.Patch("/api/bots/:botId<int>/components/:compId<int>", h.UpdComponent(app.Db, app.Redis, app.Log))

	// Set next step component
	app.Server.Post("/api/bots/:botId<int>/components/:compId<int>/next", h.SetNextStepComponent(app.Db, app.Redis, app.Log))

	// Get bot components
	app.Server.Get("/api/bots/:botId<int>/components", h.GetBotComponents(app.Db, app.Log))

	// Delete next step component
	app.Server.Delete("/api/bots/:botId<int>/components/:compId<int>/next", h.DelNextStepComponent(app.Db, app.Redis, app.Log))
}

func (app *App) regCommandsHandlers() {
	// Add command
	app.Server.Post("/api/bots/:botId<int>/components/:compId<int>/commands", h.AddCommand(app.Db, app.Redis, app.Log))

	// Delete command
	app.Server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.DelCommand(app.Db, app.Redis, app.Log))

	// update command
	app.Server.Patch("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.UpdCommand(app.Db, app.Redis, app.Log))

	// Set next step command
	app.Server.Post("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.SetNextStepCommand(app.Db, app.Redis, app.Log))

	// Delete next step command
	app.Server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.DelNextStepCommand(app.Db, app.Redis, app.Log))
}
