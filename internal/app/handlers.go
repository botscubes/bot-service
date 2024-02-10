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
	//components := app.server.Group("/api/bots/:botId<int>/groups/:groupId<int>/components")
	regComponentsHandlers(components, h)

	regComponentHandlers(component, h)
	//app.regCommandsHandlers(h)

	// custom 404 handler
	app.server.Use(handlers.NotFoundHandler)
}

// Bot handlers
func regBotsHandlers(bots fiber.Router, h *handlers.ApiHandler) {
	// Create new bot
	bots.Post("", h.NewBot)
	// Get user bots
	bots.Get("", h.GetBots)

	// Wipe bot data
	//app.server.Patch("/api/bots/:botId<int>/wipe", h.WipeBot)
}

func regBotHandlers(bot fiber.Router, h *handlers.ApiHandler) {

	// Delete bot
	bot.Delete("", h.DeleteBot)

	// Set bot token
	bot.Post("/token", h.SetBotToken)

	// Delete bot token
	bot.Delete("/token", h.DeleteBotToken)

	// Start bot
	bot.Patch("/start", h.StartBot)

	// Stop bot
	bot.Patch("/stop", h.StopBot)

}

func regGroupsHandlers(groups fiber.Router, h *handlers.ApiHandler) {

}

func regGroupHandlers(group fiber.Router, h *handlers.ApiHandler) {
	group.Post("/connections", h.AddConnetion)
}

func regComponentsHandlers(components fiber.Router, h *handlers.ApiHandler) {
	// // Adds a component to the bot structure
	// app.server.Post("/api/bots/:botId<int>/components", h.AddComponent)
	//
	// // Delete bot component
	// app.server.Delete("/api/bots/:botId<int>/components/:compId<int>", h.DelComponent)
	//
	// // Delete set of components
	// app.server.Post("/api/bots/:botId<int>/components/del", h.DelSetOfComponents)
	//
	// // update component
	// app.server.Patch("/api/bots/:botId<int>/components/:compId<int>", h.UpdComponent)
	//
	// // Set next step component
	// app.server.Post("/api/bots/:botId<int>/components/:compId<int>/next", h.SetNextStepComponent)
	//
	// Get bot components
	components.Get("", h.GetBotComponents)
	components.Post("", h.AddComponent)
	//
	// // Delete next step component
	// app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/next", h.DelNextStepComponent)
}

func regComponentHandlers(component fiber.Router, h *handlers.ApiHandler) {
	component.Delete("", h.DeleteComponent)
	component.Patch("/position", h.SetComponentPosition)
}

//func (app *App) regCommandsHandlers(h *handlers.ApiHandler) {
//	// Add command
//	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/commands", h.AddCommand)
//
//	// Delete command
//	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.DelCommand)
//
//	// update command
//	app.server.Patch("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>", h.UpdCommand)
//
//	// Set next step command
//	app.server.Post("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.SetNextStepCommand)
//
//	// Delete next step command
//	app.server.Delete("/api/bots/:botId<int>/components/:compId<int>/commands/:commandId<int>/next", h.DelNextStepCommand)
//}
