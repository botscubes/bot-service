package app

import (
	h "github.com/botscubes/bot-service/internal/api/handlers"
)

func (app *App) addHandlers() {
	app.Router.GET("/api/bot/health", h.Auth(h.Health, &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.POST("/api/bot/new", h.Auth(h.NewBot(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.POST("/api/bot/setToken", h.Auth(h.SetToken(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	// Mb change to DELETE http methon
	app.Router.POST("/api/bot/deleteToken", h.Auth(h.DeleteToken(app.Db), &app.SessionStorage, &app.Conf.JWTKey))

	app.Router.POST("/api/bot/start", h.Auth(h.StartBot(app.Db, &app.Bots, app.Server, &app.Conf.Bot), &app.SessionStorage, &app.Conf.JWTKey))
	app.Router.POST("/api/bot/stop", h.Auth(h.StopBot(app.Db, &app.Bots), &app.SessionStorage, &app.Conf.JWTKey))
}
