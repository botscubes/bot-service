package handlers

import (
	"context"
	"errors"
	"strings"

	"github.com/goccy/go-json"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/user-service/pkg/jwt"
	"github.com/botscubes/user-service/pkg/token_storage"

	fh "github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")

	ErrUserIDConvertation = errors.New("userId convertation to int64")
)

type reqHandler = fh.RequestHandler

func doJsonRes(ctx *fh.RequestCtx, code int, obj any) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fh.StatusInternalServerError)
	}
}

func Auth(h reqHandler, st *token_storage.TokenStorage, jwtKey *string, log *zap.SugaredLogger) reqHandler {
	return fh.RequestHandler(func(ctx *fh.RequestCtx) {
		const prefix = "Bearer "

		auth := ctx.Request.Header.Peek("Authorization")
		if auth == nil {
			doJsonRes(ctx, fh.StatusUnauthorized, resp.New(false, nil, e.ErrUnauthorized))
			return
		}

		token := string(auth)
		if !strings.HasPrefix(token, prefix) {
			doJsonRes(ctx, fh.StatusUnauthorized, resp.New(false, nil, e.ErrUnauthorized))
			return
		}

		token = strings.TrimPrefix(token, prefix)
		exists, err := (*st).CheckToken(context.Background(), token)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !exists {
			doJsonRes(ctx, fh.StatusUnauthorized, resp.New(false, nil, e.ErrUnauthorized))
			return
		}

		// WARN: fix error !!!
		userId, err := jwt.GetIdFromToken(token, *jwtKey)
		if err != nil {
			doJsonRes(ctx, fh.StatusUnauthorized, resp.New(false, nil, e.ErrUnauthorized))
			return
		}

		ctx.SetUserValue("userId", int64(userId))

		h(ctx)
	})
}

func Health(ctx *fh.RequestCtx) {
	_, _ = ctx.WriteString("OK")
	ctx.SetStatusCode(fh.StatusOK)
}

func PanicHandler(log *zap.SugaredLogger) func(ctx *fh.RequestCtx, err any) {
	return func(ctx *fh.RequestCtx, err any) {
		if err != nil {
			log.Errorf("API panic recovered: %v", err)
		}

		doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
	}
}
