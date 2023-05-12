package handlers

import (
	"strconv"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	"github.com/goccy/go-json"
	fh "github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type addCommandReq struct {
	Type *string `json:"type"`
	Data *string `json:"data"`
}

type addCommandRes struct {
	Id int64 `json:"id"`
}

func AddCommand(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData addCommandReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if err := model.ValidateCommand(reqData.Type, reqData.Data); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
			return
		}

		// TODO: check coomand type & data valid

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
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

func DelCommand(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
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

type setNextStepCommandReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

func SetNextStepCommand(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData setNextStepCommandReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.NextStepId == nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.MissingParam("nextStepId")))
			return
		}

		nextComponentId := reqData.NextStepId

		if *nextComponentId == compId {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.InvalidParam("nextStepId")))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
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

func DelNextStepCommand(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		compId, err := strconv.ParseInt(ctx.UserValue("compId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		commandId, err := strconv.ParseInt(ctx.UserValue("commandId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
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
