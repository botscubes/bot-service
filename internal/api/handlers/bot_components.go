package handlers

import (
	"strconv"

	"github.com/goccy/go-json"
	fh "github.com/valyala/fasthttp"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/botscubes/bot-service/pkg/log"
)

type addComponentReq struct {
	Data     *model.Data     `json:"data"`
	Commands *model.Commands `json:"commands"`
	Position *model.Point    `json:"position"`
}

type addComponentRes struct {
	Id int64 `json:"id"`
}

func AddComponent(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData addComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: AddComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// TODO: check fields limits:
		// eg. data.commands._.data max size, check commands max count
		if err := model.ValidateComponent(reqData.Data, reqData.Commands, reqData.Position); err != nil {
			log.Debug("[API: AddComponent] - [ValidateComponent];\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: AddComponent] - userId convertation to int64 error")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: AddComponent] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		component := &model.Component{
			Data: reqData.Data,
			Keyboard: &model.Keyboard{
				Buttons: [][]*int64{},
			},
			NextStepId: nil,
			IsMain:     false,
			Position:   reqData.Position,
			Status:     pgsql.StatusComponentActive,
		}

		compId, err := db.AddBotComponent(botId, component)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		for _, v := range *reqData.Commands {
			mc := &model.Command{
				Type:        v.Type,
				Data:        v.Data,
				ComponentId: &compId,
				NextStepId:  nil,
				Status:      pgsql.StatusCommandActive,
			}

			_, err := db.AddCommand(botId, mc)
			if err != nil {
				log.Error(err)
				doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
				return
			}
		}

		dataRes := &addComponentRes{
			Id: compId,
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, dataRes, nil))
	}
}

type setNextStepComponentReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextStepComponent(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextStepComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextStepComponent] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData setNextStepComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextStepComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.NextStepId == nil {
			log.Debug("[API: SetNextStepComponent] nextStepId is misssing")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.InvalidParam("nextStepId")))
			return
		}

		nextComponentId := reqData.NextStepId

		if *nextComponentId == compId {
			log.Debug("[API: SetNextStepComponent] nextStepId == stepID")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.InvalidParam("nextStepId")))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetNextStepComponent] - userId convertation to int64 error")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: SetNextStepComponent] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existInitialComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextStepComponent] - initial component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		// check bot next component exists
		existNextComp, err := db.CheckComponentExist(botId, *nextComponentId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextStepComponent] - next component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrNextComponentNotFound))
			return
		}

		if err = db.SetNextStepComponent(botId, compId, *nextComponentId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

type setNextStepCommandReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextStepCommand(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextStepCommand] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextStepCommand] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: SetNextStepCommand] - commandId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData setNextStepCommandReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: SetNextStepCommand] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.NextStepId == nil {
			log.Debug("[API: SetNextStepCommand] nextStepId is misssing")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.InvalidParam("nextStepId")))
			return
		}

		nextComponentId := reqData.NextStepId

		if *nextComponentId == compId {
			log.Debug("[API: SetNextStepCommand] nextStepId == stepID")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.InvalidParam("nextStepId")))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: SetNextStepCommand] - userId convertation to int64 error")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: SetNextStepCommand] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existInitialComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existInitialComp {
			log.Debug("[API: SetNextStepCommand] - initial component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		// Check next component exists
		existNextComp, err := db.CheckComponentExist(botId, *nextComponentId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existNextComp {
			log.Debug("[API: SetNextStepCommand] - next component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrNextComponentNotFound))
			return
		}

		// Check command exists
		existCommand, err := db.CheckCommandExist(botId, compId, commandId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existCommand {
			log.Debug("[API: SetNextStepCommand] - command not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrCommandNotFound))
			return
		}

		if err = db.SetNextStepCommand(botId, commandId, *nextComponentId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func GetBotComponents(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: GetBotComponents] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: GetBotComponents] - userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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

		components, err := db.BotComponentsForEd(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, components, nil))
	}
}

func DelNextStepComponent(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelNextStepComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelNextStepComponent] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DelNextStepComponent] - userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: DelNextStepComponent] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existComp {
			log.Debug("[API: DelNextStepComponent] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		if err = db.DelNextStepComponent(botId, compId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func DelNextStepCommand(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelNextStepCommand] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelNextStepCommand] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelNextStepCommand] - commandId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DelNextStepCommand] - userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: DelNextStepCommand] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existComp {
			log.Debug("[API: DelNextStepCommand] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		// Check command exists
		existCommand, err := db.CheckCommandExist(botId, compId, commandId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existCommand {
			log.Debug("[API: DelNextStepCommand] - command not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrCommandNotFound))
			return
		}

		if err = db.DelNextStepCommand(botId, commandId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func DelBotComponent(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelBotComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelBotComponent] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// check component is main
		if compId == config.MainComponentId {
			log.Debug("[API: DelBotComponent] - component is main;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrMainComponent))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DelBotComponent] - userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: DelBotComponent] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existComp {
			log.Debug("[API: DelBotComponent] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		if err = db.DelBotComponent(botId, compId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if err = db.DelCommandsByCompId(botId, compId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

func DelCommand(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelCommand] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelNextStepCommand] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelCommand] - commandId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DelCommand] - userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: DelCommand] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existComp {
			log.Debug("[API: DelCommand] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		// Check command exists
		existCommand, err := db.CheckCommandExist(botId, compId, commandId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existCommand {
			log.Debug("[API: DelCommand] - command not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrCommandNotFound))
			return
		}

		if err = db.DelCommand(botId, commandId); err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}

type addCommandReq struct {
	Type *string `json:"type"`
	Data *string `json:"data"`
}

type addCommandRes struct {
	Id int64 `json:"id"`
}

func AddCommand(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddCommand] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: AddCommand] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData addCommandReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: AddCommand] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if err := model.ValidateCommand(reqData.Type, reqData.Data); err != nil {
			log.Debug("[API: AddCommand] - [ValidateCommand];\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
			return
		}

		// TODO: check coomand type & data valid

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: AddCommand] - userId convertation to int64 error;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
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
			log.Debug("[API: AddCommand] - bot not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrBotNotFound))
			return
		}

		// check bot component exists
		existComp, err := db.CheckComponentExist(botId, compId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existComp {
			log.Debug("[API: DelBotComponent] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		m := &model.Command{
			Type:        reqData.Type,
			Data:        reqData.Data,
			ComponentId: &compId,
			NextStepId:  nil,
			Status:      pgsql.StatusCommandActive,
		}

		commandId, err := db.AddCommand(botId, m)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		// Invalidate component cache
		if err = r.DelComponent(botId, compId); err != nil {
			log.Error(err)
		}

		dataRes := &addCommandRes{
			Id: commandId,
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, dataRes, nil))
	}
}
