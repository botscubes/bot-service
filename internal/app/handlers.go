package app

import (
	"github.com/botscubes/bot-service/internal/api/handlers"
	m "github.com/botscubes/bot-service/internal/api/middlewares"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (app *App) regiterHandlers(h *handlers.ApiHandler) {
	app.server.Get("/api/bots/health", handlers.Health)

	// panic recover
	app.server.Use(recover.New())

	// Auth middleware
	app.server.Use(m.Auth(&app.sessionStorage, &app.conf.JWTKey, app.log))

	app.regBotsHandlers(h)
	app.regComponentsHandlers(h)
	app.regCommandsHandlers(h)

	// custom 404 handler
	app.server.Use(handlers.NotFoundHandler)
}

// Bot handlers
func (app *App) regBotsHandlers(h *handlers.ApiHandler) {
	// Create new bot
	app.server.Post("/api/bots", h.NewBot)

	// Set bot token
	app.server.Post("/api/bots/:botId<int>/token", h.SetBotToken)

	// Delete bot token
	app.server.Delete("/api/bots/:botId<int>/token", h.DeleteBotToken)

	// Start bot
	app.server.Patch("/api/bots/:botId<int>/start", h.StartBot)

	// Stop bot
	app.server.Patch("/api/bots/:botId<int>/stop", h.StopBot)

	// Get user bots
	app.server.Get("/api/bots", h.GetBots)

	// Wipe bot data
	app.server.Patch("/api/bots/:botId<int>/wipe", h.WipeBot)
}

func (app *App) regComponentsHandlers(h *handlers.ApiHandler) {
	// Adds a component to the bot structure
	app.server.Post("/api/bots/:botId<int>/components", h.AddComponent)

	// Delete bot component
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>", h.DelComponent)

	// Delete set of components
	app.server.Post("/api/bots/:botId<int>/components/del", h.DelSetOfComponents)

	// update component
	app.server.Patch("/api/bots/:botId<int>/components/:compId<int>", h.UpdComponent)

	// Set next step component
	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/next", h.SetNextStepComponent)

	// Get bot components
	app.server.Get("/api/bots/:botId<int>/components", h.GetBotComponents)

	// Delete next step component
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/next", h.DelNextStepComponent)
}

func (app *App) regCommandsHandlers(h *handlers.ApiHandler) {
	// Add command
	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/commands", h.AddCommand)

	// Delete command
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.DelCommand)

	// update command
	app.server.Patch("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.UpdCommand)

	// Set next step command
	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.SetNextStepCommand)

	// Delete next step command
	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.DelNextStepCommand)
}
