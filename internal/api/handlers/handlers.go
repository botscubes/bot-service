package handlers

import (
	"encoding/json"
	"strings"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/botscubes/user-service/pkg/jwt"
	"github.com/botscubes/user-service/pkg/token_storage"
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

func auth(h fasthttp.RequestHandler, st token_storage.TokenStorage, jwtKey *string) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		const prefix = "Bearer "

		auth := ctx.Request.Header.Peek("Authorization")
		if auth == nil {
			doJsonRes(ctx, fasthttp.StatusUnauthorized, resp.New(false, nil, errors.ErrUnauthorized))
			return
		}
		token := string(auth)
		if !strings.HasPrefix(token, prefix) {
			doJsonRes(ctx, fasthttp.StatusUnauthorized, resp.New(false, nil, errors.ErrUnauthorized))
			return
		}

		token = strings.TrimPrefix(token, prefix)
		if token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.Vl16d9RIxtWDeGXgh3cdK-KRvesGhjr96qcYqDncj8k" {
			doJsonRes(ctx, fasthttp.StatusUnauthorized, resp.New(false, nil, errors.ErrUnauthorized))
			return
		}

		exists, err := st.CheckToken(token)
		if err != nil {
			log.Error("[API: auth middleware] [CheckToken]\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !exists {
			doJsonRes(ctx, fasthttp.StatusUnauthorized, resp.New(false, nil, errors.ErrUnauthorized))
			return
		}

		// WARN: fix error !!!
		id, err := jwt.GetIdFromToken(token, *jwtKey)
		if err != nil {
			doJsonRes(ctx, fasthttp.StatusUnauthorized, resp.New(false, nil, errors.ErrUnauthorized))
			return
		}

		log.Info(id)

		h(ctx)
	})
}

func health(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("OK")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func AddHandlers(r *fastRouter.Router, db *pgsql.Db, bots *map[string]*bot.TBot, server *telego.MultiBotWebhookServer, conf *config.BotConfig, st token_storage.TokenStorage, jwtKey *string) {
	r.GET("/api/bot/health", auth(health, st, jwtKey))

	r.POST("/api/bot/new", newBot(db))
	r.POST("/api/bot/setToken", setToken(db))

	// Mb change to DELETE http methon
	r.POST("/api/bot/deleteToken", deleteToken(db))

	r.POST("/api/bot/start", startBot(db, bots, server, conf))
	r.POST("/api/bot/stop", stopBot(db, bots))
}
