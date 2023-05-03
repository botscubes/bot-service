package handlers

import (
	"strconv"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/pgtype"
	fh "github.com/valyala/fasthttp"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
)

type component struct {
	Id         int64          `json:"id"`
	Data       *componentData `json:"data"`
	Keyboard   *keyboard      `json:"keyboard"`
	Commands   *[]*command    `json:"commands"`
	NextStepId *int64         `json:"nextStepId"`
	IsStart    bool           `json:"isStart"`
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
	Type    *string         `json:"type"`
	Content *[]*dataContent `json:"content"`
}

type dataContent struct {
	Text *string `json:"text,omitempty"`
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
	return func(ctx *fh.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddBotComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData addBotComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: AddBotComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// TODO: check fields limits:
		// eg. data.commands._.data max size

		if err := validateAddBotComponent(&reqData); err != nil {
			log.Debug(err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidParams))
			return
		}

		px := reqData.Position.X
		py := reqData.Position.Y

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: AddBotComponent] - get userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: AddBotComponent] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		status := 0

		m := &model.Component{
			Data: &model.Data{
				Type:    reqData.Data.Type,
				Content: *dataContentsMod(reqData.Data.Content),
			},
			Keyboard: &model.Keyboard{
				Buttons: [][]*int64{},
			},
			NextStepId: nil,
			IsStart:    false,
			Position: &pgtype.Point{
				P:      pgtype.Vec2{X: *px, Y: *py},
				Status: pgtype.Present,
			},
			Status: status,
		}

		compId, err := db.AddBotComponent(botId, m)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
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
					log.Error(err)
					doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
					return
				}
			}
		}

		dataRes := &addBotComponentRes{
			Id: compId,
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, dataRes, nil))
	}
}

type setNextForComponentReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextForComponent(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForComponent] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData setNextForComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextForComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.NextStepId == nil {
			log.Debug("[API: SetNextForComponent] nextStepId is misssing")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidParams))
			return
		}

		nextComponentId := reqData.NextStepId
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetNextForComponent] - get userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: SetNextForComponent] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existInitialComp, err := db.CheckBotComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextForComponent] - initial component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		// check bot next component exists
		existNextComp, err := db.CheckBotComponentExist(botId, *nextComponentId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextForComponent] - next component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrNextComponentNotFound))
			return
		}

		if err = db.SetNextStepForComponent(botId, compId, *nextComponentId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

type SetNextForCommandReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextForCommand(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextForCommand] - commandId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData SetNextForCommandReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextForCommand] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.NextStepId == nil {
			log.Debug("[API: SetNextForCommand] nextStepId is misssing")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidParams))
			return
		}

		nextComponentId := reqData.NextStepId
		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetNextForCommand] - get userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: SetNextForCommand] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// Check initial component exists
		existInitialComp, err := db.CheckBotComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextForCommand] - initial component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		// Check next component exists
		existNextComp, err := db.CheckBotComponentExist(botId, *nextComponentId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextForCommand] - next component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrNextComponentNotFound))
			return
		}

		// Check command exists
		existCommand, err := db.CheckBotCommandExist(botId, compId, commandId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existCommand {
			log.Debug("[API: SetNextForCommand] - command not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrCommandNotFound))
			return
		}

		if err = db.SetNextStepForCommand(botId, compId, *nextComponentId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

type getBotCompsRes []*component

func GetBotComponents(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		var err error

		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: GetBotComponents] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: GetBotComponents] - get userId convertation to int64 error;\n")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: GetBotComponents] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		components, err := db.GetBotComponents(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, componentsRes(components), nil))
	}
}
