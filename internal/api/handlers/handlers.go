package handlers

import (
	"encoding/json"

	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	fastRouter "github.com/fasthttp/router"
	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
)

// TODO: check user access
// Add check all id's for > 0

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func doJsonRes(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func health(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("OK")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func AddHandlers(r *fastRouter.Router, db *pgsql.Db, bots *map[string]*bot.TBot, server *telego.MultiBotWebhookServer, conf *config.BotConfig) {
	r.GET("/api/bot/health", health)

	r.POST("/api/bot/new", newBot(db))
	r.POST("/api/bot/setToken", setToken(db))

	// Mb change to DELETE http methon
	r.POST("/api/bot/deleteToken", deleteToken(db))

	r.POST("/api/bot/start", startBot(db, bots, server, conf))
}
