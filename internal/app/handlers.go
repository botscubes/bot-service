package app

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unicode/utf8"

	errors "github.com/botscubes/bot-service/internal/api/errors"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/valyala/fasthttp"
)

type newBotReq struct {
	UserId *json.Number `json:"user_id"`
	Title  *string      `json:"title,omitempty"`
}

type newBotRes struct {
	Id int64 `json:"id"`
}

// type newBotJson struct {
// 	Token string `json:"token"`
// }

const (
	tokenRegexp = `^\d{9,10}:[\w-]{35}$` //nolint:gosec
	maxTitleLen = 50
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func validateToken(token string) bool {
	reg := regexp.MustCompile(tokenRegexp)
	return reg.MatchString(token)
}

func doJsonRes(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func newBotHandler(ctx *fasthttp.RequestCtx) {
	log.Debug("[API: newBotHandler] - Start")

	var err error = nil

	var data newBotReq
	err = json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Debug("[API: newBotHandler] - Serialisation error;\n %s", err)
		doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
		return
	}

	if data.UserId == nil {
		log.Debug("[API: newBotHandler] user_id is misssing")
		doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidParams))
		return
	}

	user_id, err := data.UserId.Int64()
	if err != nil {
		log.Debug("[API: newBotHandler] - (user_id) json.Number convertation to int64 error;\n %s", err)
		doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
		return
	}

	title := data.Title
	token := ""
	status := 0

	if title == nil || *title == "" {
		log.Debug("[API: newBotHandler] - title is misssing")
		doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidParams))
		return
	}

	if utf8.RuneCountInString(*title) > maxTitleLen {
		log.Debug("[API: newBotHandler] - title len > ", maxTitleLen)
		doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidTitleLength))
		return
	}

	botId, err := app.db.AddBot(user_id, &token, title, status)
	if err != nil {
		log.Debug("[API: newBotHandler] - db AppBot error;\n %s", err)
		doJsonRes(ctx, fasthttp.StatusOK, resp.New(false, nil, errors.ErrInvalidRequest))
		return
	}

	log.Info(botId)

	dataRes := &newBotRes{
		Id: botId,
	}

	doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
}

// Edit to StartBot
// func newBotHandler(ctx *fasthttp.RequestCtx) {
// 	// TODO: Check token exist (from db)

// 	var err error = nil
// 	log.Debug("[API: newBotHandler] - NEW TOKEN")

// 	var data newBotJson
// 	err = json.Unmarshal(ctx.PostBody(), &data)
// 	if err != nil {
// 		log.Errorf("[API: newBotHandler] - Serialisation error;\n %s", err)
// 		doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrInvalidRequest)
// return
// 	}

// 	token := data.Token

// 	if !validateToken(token) {
// 		log.Debug("[API: newBotHandler] Incorrect token")
// 		doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrIncorrectTokenFormat)
// 		return
// 	}

// 	// TODO: Own token health check to get a specific error
// 	nbot, err := bot.NewBot(token)
// 	if err != nil {
// 		log.Debug("[API: newBotHandler] ", err)
// 		doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrInvalidToken)
// 		return
// 	}

// 	if _, ok := app.bots[token]; ok {
// 		log.Debug("[API: newBotHandler] Token already exist in bots map")
// 		doJsonRes(ctx, fasthttp.StatusOK, &errors.ErrTokenExistInSystem)
// 		return
// 	}

// 	app.bots[token] = new(bot.TBot)
// 	app.bots[token].Bot = nbot

// 	doJsonRes(ctx, fasthttp.StatusOK, &errors.Success)
// }

func healthHandler(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("OK")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func startBotHandler(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString(fmt.Sprintf("Started: %s", ctx.UserValue("botid")))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
