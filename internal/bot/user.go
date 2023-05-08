package bot

import (
	"errors"

	"github.com/botscubes/bot-service/internal/database/pgsql"
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

	if err.Error() != "not found" {
		log.Error(err)
	}

	// userStep not found in cache, try get from db
	exist, err := btx.Db.CheckUserExistByTgId(btx.Id, from.ID)
	if err != nil {
		return 0, err
	}

	if exist {
		stepID, err = btx.Db.UserStepByTgId(btx.Id, from.ID)
		if err != nil {
			return 0, err
		}

		if err = btx.Rdb.SetUserStep(btx.Id, from.ID, stepID); err != nil {
			log.Error(err)
		}

		return stepID, nil
	}

	return 0, errors.New("user not found")
}

func (btx *TBot) addUser(from *telego.User) error {
	user := &model.User{
		TgId:      from.ID,
		FirstName: &from.FirstName,
		LastName:  &from.LastName,
		Username:  &from.Username,
		StepID: model.StepID{
			StepId: 1,
		},
		Status: pgsql.StatusUserActive,
	}

	_, err := btx.Db.AddUser(btx.Id, user)
	if err != nil {
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
	if err != nil && err.Error() != "not found" {
		log.Error(err)
	}

	if component != nil {
		return component, nil
	}

	// component not found in cache, try get from db

	exist, err := btx.Db.CheckComponentExist(btx.Id, stepID)
	if err != nil {
		return nil, err
	}

	if exist {
		component, err = btx.Db.ComponentForBot(btx.Id, btx.Id)
		if err != nil {
			return nil, err
		}
		// TODO: push to cache
		return component, nil
	}

	return nil, errors.New("not found")
}
