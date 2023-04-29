package handlers

import (
	"strconv"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/pgtype"

	"github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/valyala/fasthttp"
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
	Id      *int64  `json:"id,omitempty"`
	Type    *string `json:"type"`
	Data    *string `json:"data"`
	Next_id *int64  `json:"next_id,omitempty"`
}

type addComponentRes struct {
	Id int64 `json:"id"`
}

func AddComponent(db *pgsql.Db) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		bot_id, err := strconv.ParseInt(ctx.UserValue("bot_id").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddComponent] - bot_id param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		var reqData addComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: AddComponent] - Serialisation error;\n", err)
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

		user_id := ctx.UserValue("user_id").(int64)

		existBot, err := db.CheckBotExist(user_id, bot_id)
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

		compId, err := db.AddComponent(bot_id, m)
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

				_, err := db.AddCommand(bot_id, mc)
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
	NextId *int64 `json:"next_id"`
}

func SetNextForComponent(db *pgsql.Db) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error = nil

		bot_id, err := strconv.ParseInt(ctx.UserValue("bot_id").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - bot_id param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		comp_id, err := strconv.ParseInt(ctx.UserValue("comp_id").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - comp_id param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		var reqData setNextForComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextForComponent] - Serialisation error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if reqData.NextId == nil {
			log.Debug("[API: SetNextForComponent] next_id is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		nextId := reqData.NextId
		user_id := ctx.UserValue("user_id").(int64)

		existBot, err := db.CheckBotExist(user_id, bot_id)
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

		existInitialComp, err := db.CheckComponentExist(bot_id, comp_id)
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

		existNextComp, err := db.CheckComponentExist(bot_id, *nextId)
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

		if err = db.SetNextIdForComponent(bot_id, comp_id, *nextId); err != nil {
			log.Debug("[API: SetNextForComponent] - [db: SetNextIdForComponent] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}
