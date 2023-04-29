package app

import (
	h "github.com/botscubes/bot-service/internal/api/handlers"
)

func (app *App) addHandlers() {
	// TODO: Correct methods (create, delete, ...)

	app.Router.GET("/api/bots/health", h.Auth(h.Health, &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.POST("/api/bots/new", h.Auth(h.NewBot(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.POST("/api/bots/{bot_id}/setToken", h.Auth(h.SetToken(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.DELETE("/api/bots/{bot_id}/deleteToken", h.Auth(h.DeleteToken(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.PATCH("/api/bots/{bot_id}/start", h.Auth(h.StartBot(app.Db, &app.Bots, app.Server, &app.Conf.Bot), &app.SessionStorage, &app.Conf.JWTKey))
	app.Router.PATCH("/api/bots/{bot_id}/stop", h.Auth(h.StopBot(app.Db, &app.Bots), &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.POST("/api/bots/{bot_id}/component/add", h.Auth(h.AddComponent(app.Db), &app.SessionStorage, &app.Conf.JWTKey))
}
