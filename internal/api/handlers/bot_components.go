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

// TODO:
// add commands ids to components keyboard
// select command from components keyboard

type component struct {
	Id         int64          `json:"id"`
	Data       *componentData `json:"data"`
	Keyboard   *keyboard      `json:"keyboard"`
	Commands   *[]*command    `json:"commands"`
	NextStepId *int64         `json:"nextStepId"`
	Position   *point         `json:"position"`
}

type command struct {
	Id          *int64  `json:"id,omitempty"`
	Type        *string `json:"type"`
	Data        *string `json:"data"`
	ComponentId *int64  `json:"componentId"`
	NextStepId  *int64  `json:"nextStepId"`
}

type point struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
}

type componentData struct {
	Type    *string      `json:"type"`
	Content *dataContent `json:"content"`
}

type dataContent struct {
	Text *string `json:"text"`
}

type keyboard struct {
	Buttons [][]*int64 `json:"buttons"`
}

// all field reqired
// comands: [] | []command
type addBotComponentReq struct {
	Data     *componentData `json:"data"`
	Commands []*command     `json:"commands"`
	Position *point         `json:"position"`
}

type addBotComponentRes struct {
	Id int64 `json:"id"`
}

func AddBotComponent(db *pgsql.Db) reqHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddBotComponent] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		var reqData addBotComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: AddBotComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		// TODO: check fields limits:
		// eg. data.commands._.data max size

		if err := validateBotComponent(&reqData); err != nil {
			log.Debug(err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		px := reqData.Position.X
		if px != nil {
			log.Debug("[API: AddBotComponent] Position.X is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		py := reqData.Position.Y
		if py != nil {
			log.Debug("[API: AddBotComponent] Position.Y is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: AddBotComponent] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: AddBotComponent] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: AddBotComponent] - bot not found")
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
			NextStepId: nil,
			Position: &pgtype.Point{
				P:      pgtype.Vec2{X: *px, Y: *py},
				Status: pgtype.Present,
			},
			Status: status,
		}

		compId, err := db.AddBotComponent(botId, m)
		if err != nil {
			log.Error("[API: AddBotComponent] - [db: AddBotComponent] error;\n", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		// TODO: check commands max count
		if reqData.Commands != nil {
			commandStatus := 0
			for _, v := range reqData.Commands {
				mc := &model.Command{
					Type:        v.Type,
					Data:        v.Data,
					ComponentId: &compId,
					NextStepId:  nil,
					Status:      commandStatus,
				}

				_, err := db.AddBotCommand(botId, mc)
				if err != nil {
					log.Error("[API: AddBotComponent] - [db: AddBotCommand] error;\n", err)
					doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
					return
				}
			}
		}

		dataRes := &addBotComponentRes{
			Id: compId,
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, dataRes, nil))
	}
}

type setNextForComponentReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextForComponent(db *pgsql.Db) reqHandler {
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

		if reqData.NextStepId == nil {
			log.Debug("[API: SetNextForComponent] nextStepId is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		nextComponentId := reqData.NextStepId
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

		existInitialComp, err := db.CheckBotComponentExist(botId, compId)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - [db: CheckBotComponentExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextForComponent] - initial component not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrComponentNotFound))
			return
		}

		existNextComp, err := db.CheckBotComponentExist(botId, *nextComponentId)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - [db: CheckBotComponentExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextForComponent] - next component not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrNextComponentNotFound))
			return
		}

		if err = db.SetNextStepForComponent(botId, compId, *nextComponentId); err != nil {
			log.Debug("[API: SetNextForComponent] - [db: SetNextStepForComponent] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}

type SetNextForCommandReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextForCommand(db *pgsql.Db) reqHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - compId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - commandId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		var reqData SetNextForCommandReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextForCommand] - Serialization error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		if reqData.NextStepId == nil {
			log.Debug("[API: SetNextForCommand] nextStepId is misssing")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidParams))
			return
		}

		nextComponentId := reqData.NextStepId
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetNextForCommand] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: SetNextForCommand] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		// Check initial component exists
		existInitialComp, err := db.CheckBotComponentExist(botId, compId)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - [db: CheckBotComponentExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextForCommand] - initial component not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrComponentNotFound))
			return
		}

		// Check next component exists
		existNextComp, err := db.CheckBotComponentExist(botId, *nextComponentId)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - [db: CheckBotComponentExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextForCommand] - next component not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrNextComponentNotFound))
			return
		}

		// Check command exists
		existCommand, err := db.CheckBotCommandExist(botId, compId, commandId)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - [db: CheckBotCommandExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existCommand {
			log.Debug("[API: SetNextForCommand] - command not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrCommandNotFound))
			return
		}

		if err = db.SetNextStepForCommand(botId, compId, *nextComponentId); err != nil {
			log.Debug("[API: SetNextForCommand] - [db: SetNextStepForCommand] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, nil, nil))
	}
}

type botFullCompsRes []*component

func GetBotComponents(db *pgsql.Db) reqHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: GetBotComponents] - botId param error;\n", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: GetBotComponents] - get userId convertation to int64 error;", err)
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrInvalidRequest))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Debug("[API: GetBotComponents] - [db: CheckBotExist] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: GetBotComponents] - bot not found")
			doJsonRes(ctx, fasthttp.StatusBadRequest, resp.New(false, nil, errors.ErrBotNotFound))
			return
		}

		components, err := db.GetBotFullComponents(botId)
		if err != nil {
			log.Debug("[API: GetBotComponents] - [db: GetBotFullComponents] error;", err)
			doJsonRes(ctx, fasthttp.StatusInternalServerError, resp.New(false, nil, errors.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fasthttp.StatusOK, resp.New(true, botFullComponentsRes(components), nil))
	}
}
