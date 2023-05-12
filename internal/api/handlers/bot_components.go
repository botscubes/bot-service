package handlers

import (
	"strconv"

	"github.com/goccy/go-json"
	fh "github.com/valyala/fasthttp"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	resp "github.com/botscubes/bot-service/pkg/api_response"
)

type addComponentReq struct {
	Data     *model.Data     `json:"data"`
	Commands *model.Commands `json:"commands"`
	Position *model.Point    `json:"position"`
}

type addComponentRes struct {
	Id int64 `json:"id"`
}

func AddComponent(db *pgsql.Db, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
		if err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		var reqData addComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		// TODO: check fields limits:
		// eg. data.commands._.data max size, check commands max count
		if err := model.ValidateComponent(reqData.Data, reqData.Commands, reqData.Position); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
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

func SetNextStepComponent(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
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

		var reqData setNextStepComponentReq
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

		// check bot next component exists
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

func GetBotComponents(db *pgsql.Db, log *zap.SugaredLogger) reqHandler {
	return func(ctx *fh.RequestCtx) {
		botId, err := strconv.ParseInt(ctx.UserValue("botId").(string), 10, 64)
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

		components, err := db.ComponentsForEd(botId)
		if err != nil {
			log.Error(err)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		doJsonRes(ctx, fh.StatusOK, resp.New(true, components, nil))
	}
}

func DelNextStepComponent(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
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

func DelComponent(db *pgsql.Db, r *rdb.Rdb, log *zap.SugaredLogger) reqHandler {
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

		// check component is main
		if compId == config.MainComponentId {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrMainComponent))
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

func UpdComponent(db *pgsql.Db, log *zap.SugaredLogger) reqHandler {
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

		var reqData updComponentReq
		if err = json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		userId, ok := ctx.UserValue("userId").(int64)
		if !ok {
			log.Error(ErrUserIDConvertation)
			doJsonRes(ctx, fh.StatusInternalServerError, resp.New(false, nil, e.ErrInternalServer))
			return
		}

		if reqData.Data == nil && reqData.Position == nil {
			doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, e.ErrInvalidRequest))
			return
		}

		if reqData.Data != nil {
			if err := reqData.Data.Validate(); err != nil {
				doJsonRes(ctx, fh.StatusBadRequest, resp.New(false, nil, err))
				return
			}
		}

		if reqData.Position != nil {
			if err := reqData.Position.Validate(); err != nil {
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
