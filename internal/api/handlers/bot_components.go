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

		compId, err := db.AddComponent(botId, component)
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
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.MissingParam("nextStepId")))
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

		components, err := db.ComponentsForEd(botId)
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

func DelComponent(db *pgsql.Db, r *rdb.Rdb) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: DelComponent] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// check component is main
		if compId == config.MainComponentId {
			log.Debug("[API: DelComponent] - component is main;")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrMainComponent))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: DelComponent] - userId convertation to int64 error;")
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
			log.Debug("[API: DelComponent] - bot not found")
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
			log.Debug("[API: DelComponent] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		if err = db.DelComponent(botId, compId); err != nil {
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

type updComponentReq struct {
	Data     *model.Data  `json:"data"`
	Position *model.Point `json:"position"`
}

func UpdComponent(db *pgsql.Db) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: UpdComponent] - botId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			log.Debug("[API: UpdComponent] - compId param error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData updComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			log.Debug("[API: UpdComponent] - Serialization error;\n", err)
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Debug("[API: UpdComponent] - userId convertation to int64 error")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if reqData.Data == nil && reqData.Position == nil {
			log.Debug("[API: UpdComponent] - [request fields not found]")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.Data != nil {
			if err := reqData.Data.Validate(); err != nil {
				log.Debug("[API: UpdComponent] - [Validate data];\n", err)
				doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
				return
			}
		}

		if reqData.Position != nil {
			if err := reqData.Position.Validate(); err != nil {
				log.Debug("[API: UpdComponent] - [Validate position];\n", err)
				doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
				return
			}
		}

		// check bot exists
		existBot, err := db.CheckBotExist(userId, botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if !existBot {
			log.Debug("[API: UpdComponent] - bot not found")
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
			log.Debug("[API: UpdComponent] - component not found")
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrComponentNotFound))
			return
		}

		if reqData.Position != nil {
			err = db.UpdComponentPosition(botId, compId, reqData.Position)
			if err != nil {
				log.Error(err)
				doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
				return
			}
		}

		if reqData.Data != nil {
			err = db.UpdComponentData(botId, compId, reqData.Data)
			if err != nil {
				log.Error(err)
				doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
				return
			}
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, nil, nil))
	}
}
