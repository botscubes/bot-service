package app

import (
	h "github.com/botscubes/bot-service/internal/api/handlers"
	m "github.com/botscubes/bot-service/internal/api/middlewares"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (app *App) regiterHandlers() {
	// Middlewares
	app.server.Use(recover.New())
	app.server.Use(m.Auth(&app.sessionStorage, &app.conf.JWTKey, app.log))

	app.server.Get("/api/bots/health", h.Health)

	app.regBotsHandlers()
	app.regComponentsHandlers()
	app.regCommandsHandlers()

	// custom 404 handler
	app.server.Use(h.NotFoundHandler)
}

// Bot handlers
func (app *App) regBotsHandlers() {
	// Create new bot
	app.server.Post("/api/bots", h.NewBot(app.db, app.log))

	// Set bot token
	app.server.Post("/api/bots/:botId<int>/token", h.SetBotToken(app.db, app.log))

	// Delete bot token
	app.server.Delete("/api/bots/:botId<int>/token", h.DeleteBotToken(app.db, app.log))

	// Start bot
	app.server.Patch("/api/bots/:botId<int>/start", h.StartBot(app.db, app.botService, app.log, app.nc))

	// Stop bot
	app.server.Patch("/api/bots/:botId<int>/stop", h.StopBot(app.db, app.botService, app.log, app.nc))

	// Get user bots
	app.server.Get("/api/bots", h.GetBots(app.db, app.log))

	// Wipe bot data
	app.server.Patch("/api/bots/:botId<int>/wipe", h.WipeBot(app.db, app.redis, app.botService, app.log))
}

func (app *App) regComponentsHandlers() {
	// Adds a component to the bot structure
	app.server.Post("/api/bots/:botId<int>/components", h.AddComponent(app.db, app.log))

	// Delete bot component
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>", h.DelComponent(app.db, app.redis, app.log))

	// Delete set of components
	app.server.Post("/api/bots/:botId<int>/components/del", h.DelSetOfComponents(app.db, app.redis, app.log))

	// update component
	app.server.Patch("/api/bots/:botId<int>/components/:compId<int>", h.UpdComponent(app.db, app.redis, app.log))

	// Set next step component
	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/next", h.SetNextStepComponent(app.db, app.redis, app.log))

	// Get bot components
	app.server.Get("/api/bots/:botId<int>/components", h.GetBotComponents(app.db, app.log))

	// Delete next step component
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/next", h.DelNextStepComponent(app.db, app.redis, app.log))
}

func (app *App) regCommandsHandlers() {
	// Add command
	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/commands", h.AddCommand(app.db, app.redis, app.log))

	// Delete command
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.DelCommand(app.db, app.redis, app.log))

	// update command
	app.server.Patch("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.UpdCommand(app.db, app.redis, app.log))

	// Set next step command
	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.SetNextStepCommand(app.db, app.redis, app.log))

	// Delete next step command
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.DelNextStepCommand(app.db, app.redis, app.log))
}
