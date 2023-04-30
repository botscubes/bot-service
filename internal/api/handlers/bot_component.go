package handlers

import (
	"strconv"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/pgtype"
	"github.com/valyala/fasthttp"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
)

// all field reqired
// comands: [] | []command
type addComponentReq struct {
	Data     *data      `json:"data"`
	Commands []*command `json:"commands"`
	Position *point     `json:"position"`
}

type point struct {
	X *json.Number `json:"x"`
	Y *json.Number `json:"y"`
}

type data struct {
	Type    *string  `json:"type"`
	Content *content `json:"content"`
}

type content struct {
	Text *string `json:"text"`
}

type command struct {
	Id     *int64  `json:"id,omitempty"`
	Type   *string `json:"type"`
	Data   *string `json:"data"`
	NextId *int64  `json:"nextId,omitempty"`
}

type addComponentRes struct {
	Id int64 `json:"id"`
}

func AddComponent(db *pgsql.Db) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddComponent] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		var reqData addComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: AddComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		// TODO: check fields limits:
		// eg. data.commands._.data max size

		if err := validateComponent(&reqData); err != nil {
			log.Debug(err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		px, err := reqData.Position.X.Float64()
		if err != nil {
			log.Debug("[API: AddComponent] - (position.X) json.Number convertation to Float64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		py, err := reqData.Position.Y.Float64()
		if err != nil {
			log.Debug("[API: AddComponent] - (position.Y) json.Number convertation to Float64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: AddComponent] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: AddComponent] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: AddComponent] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		status := 0

		m := &model.Component{
			Data: &model.Data{
				Type: reqData.Data.Type,
				Content: &model.Content{
					Text: reqData.Data.Content.Text,
				},
			},
			Keyboard: &model.Keyboard{
				Buttons: [][]*int64{},
			},
			NextId: nil,
			Position: &pgtype.Point{
				P:      pgtype.Vec2{X: px, Y: py},
				Status: pgtype.Present,
			},
			Status: status,
		}

		compId, err := db.AddComponent(botId, m)
		if err != nil {
			log.Error("[API: AddComponent] - [db: AddComponent] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		// TODO: check commands max count
		if reqData.Commands != nil {
			for _, v := range reqData.Commands {
				mc := &model.Command{
					Type:   v.Type,
					Data:   v.Data,
					NextId: nil,
				}

				_, err := db.AddCommand(botId, mc)
				if err != nil {
					log.Error("[API: AddComponent] - [db: AddCommand] error;\n", err)
					doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
					return
				}
			}
		}

		dataRes := &addComponentRes{
			Id: compId,
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}

type setNextForComponentReq struct {
	NextId *int64 `json:"nextId"`
}

func SetNextForComponent(db *pgsql.Db) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - compId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		var reqData setNextForComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextForComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if reqData.NextId == nil {
			log.Debug("[API: SetNextForComponent] nextId is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		nextId := reqData.NextId
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetNextForComponent] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: SetNextForComponent] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		existInitialComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - [db: CheckComponentExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextForComponent] - initial component not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrComponentNotFound))
			return
		}

		existNextComp, err := db.CheckComponentExist(botId, *nextId)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - [db: CheckComponentExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextForComponent] - next component not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrNextComponentNotFound))
			return
		}

		if err = db.SetNextIdForComponent(botId, compId, *nextId); err != nil {
			log.Debug("[API: SetNextForComponent] - [db: SetNextIdForComponent] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}
