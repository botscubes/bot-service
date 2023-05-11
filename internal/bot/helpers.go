package bot

import (
	"errors"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
)

func (btx *TBot) getUserStep(from *telego.User) (int64, error) {
	// try get userStep from cache
	stepID, err := btx.Rdb.GetUserStep(btx.Id, from.ID)
	if err == nil {
		return stepID, nil
	}

	if !errors.Is(err, rdb.ErrNotFound) {
		log.Error(err)
	}

	// userStep not found in cache, try get from db
	stepID, err = btx.Db.UserStepByTgId(btx.Id, from.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	if err = btx.Rdb.SetUserStep(btx.Id, from.ID, stepID); err != nil {
		log.Error(err)
	}

	return stepID, nil
}

func (btx *TBot) addUser(from *telego.User) error {
	user := &model.User{
		TgId:      from.ID,
		FirstName: &from.FirstName,
		LastName:  &from.LastName,
		Username:  &from.Username,
		StepID: model.StepID{
			StepId: config.MainComponentId,
		},
		Status: pgsql.StatusUserActive,
	}

	_, err := btx.Db.AddUser(btx.Id, user)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := btx.Rdb.SetUserStep(btx.Id, from.ID, user.StepId); err != nil {
		log.Error(err)
	}

	return nil
}

func (btx *TBot) getComponent(stepID int64) (*model.Component, error) {
	// try get component from cache
	component, err := btx.Rdb.GetComponent(btx.Id, stepID)
	if err == nil {
		return component, nil
	}

	if !errors.Is(err, rdb.ErrNotFound) {
		log.Error(err)
	}

	// check bot components exists in cache
	ex, err := btx.Rdb.CheckComponentsExist(btx.Id)
	if err != nil {
		log.Error(err)
	}

	// components not found in cache
	if err == nil && ex == 0 {
		// get all components from db
		components, err := btx.Db.ComponentsForBot(btx.Id)
		if err != nil {
			log.Error(err)
		}

		if err := btx.Rdb.SetComponents(btx.Id, components); err != nil {
			log.Error(err)
		}

		component, err := btx.Rdb.GetComponent(btx.Id, stepID)
		if err == nil {
			return component, nil
		}

		if !errors.Is(err, rdb.ErrNotFound) {
			log.Error(err)
		}
	}

	// component not found in cache, try get from db
	exist, err := btx.Db.CheckComponentExist(btx.Id, stepID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if exist {
		component, err = btx.Db.ComponentForBot(btx.Id, stepID)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if err = btx.Rdb.SetComponent(btx.Id, component); err != nil {
			log.Error(err)
		}

		return component, nil
	}

	return nil, ErrNotFound
}

func (btx *TBot) setUserStep(userId int64, stepID int64) {
	if err := btx.Db.SetUserStepByTgId(btx.Id, userId, stepID); err != nil {
		log.Error(err)
	}
}
