package handlers

import (
	"context"
	"strings"

	"github.com/goccy/go-json"

	"github.com/botscubes/bot-service/internal/api/errors"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/botscubes/user-service/pkg/jwt"
	"github.com/botscubes/user-service/pkg/token_storage"

	"github.com/valyala/fasthttp"
)

// TODO: check user access
// Add check all id's for > 0

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func doJsonRes(ctx *fasthttp.RequestCtx, code int, obj any) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func Auth(h fasthttp.RequestHandler, st *token_storage.TokenStorage, jwtKey *string) fasthttp.RequestHandler {
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
		exists, err := (*st).CheckToken(context.Background(), token)
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
		userId, err := jwt.GetIdFromToken(token, *jwtKey)
		if err != nil {
			doJsonRes(ctx, fasthttp.StatusUnauthorized, resp.New(false, nil, errors.ErrUnauthorized))
			return
		}

		ctx.SetUserValue("userId", int64(userId))

		h(ctx)
	})
}

func Health(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("OK")
	ctx.SetStatusCode(fasthttp.StatusOK)
}
