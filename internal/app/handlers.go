package app

import (
	"github.com/botscubes/bot-service/internal/api/handlers"
	m "github.com/botscubes/bot-service/internal/api/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (app *App) regiterHandlers(h *handlers.ApiHandler) {
	app.server.Get("/api/bots/health", handlers.Health)

	// panic recover
	app.server.Use(recover.New())

	// Auth middleware
	app.server.Use(m.Auth(&app.sessionStorage, &app.conf.JWTKey, app.log))

	api := app.server.Group("/api")
	bots := api.Group("/bots")
	bot := bots.Group("/:botId<int>", m.GetBotMiddleware(app.db, app.log))
	groups := bot.Group("/groups")
	group := groups.Group("/:groupId<int>", m.GetGroupMiddleware(app.db, app.log))
	components := group.Group("/components")
	component := components.Group("/:componentId<int>", m.GetComponentMiddleware(app.db, app.log))

	regBotsHandlers(bots, h)
	regBotHandlers(bot, h)
	regGroupHandlers(group, h)
	regComponentsHandlers(components, h)

	regComponentHandlers(component, h)

	// custom 404 handler
	app.server.Use(handlers.NotFoundHandler)
}

// Bot handlers
func regBotsHandlers(bots fiber.Router, h *handlers.ApiHandler) {
	// Create new bot
	bots.Post("", h.NewBot)
	// Get user bots
	bots.Get("", h.GetBots)

}

func regBotHandlers(bot fiber.Router, h *handlers.ApiHandler) {

	// Delete bot
	bot.Delete("", h.DeleteBot)

	// Set bot token
	bot.Patch("/token", h.SetBotToken)

	// Delete bot token
	bot.Delete("/token", h.DeleteBotToken)

	bot.Get("/token", h.GetBotToken)
	// Start bot
	bot.Patch("/start", h.StartBot)

	// Stop bot
	bot.Patch("/stop", h.StopBot)

	bot.Get("/status", h.GetBotStatus)
}

func regGroupsHandlers(groups fiber.Router, h *handlers.ApiHandler) {

}

func regGroupHandlers(group fiber.Router, h *handlers.ApiHandler) {
	group.Post("/connections", h.AddConnetion)
	group.Delete("/connections", h.DeleteConnection)
}

func regComponentsHandlers(components fiber.Router, h *handlers.ApiHandler) {

	// Get bot components
	components.Get("", h.GetBotComponents)
	components.Post("", h.AddComponent)

}

func regComponentHandlers(component fiber.Router, h *handlers.ApiHandler) {
	component.Delete("", h.DeleteComponent)
	component.Patch("/position", h.SetComponentPosition)
	component.Patch("/data", h.UpdateComponentData)
}
