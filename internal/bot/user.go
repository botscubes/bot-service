package bot

import (
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
)

func (btx *TBot) getUser(from *telego.User) (*model.User, error) {
	// try get user from cache
	user, err := btx.Rdb.GetUser(btx.Id, from.ID)
	if err != nil {
		log.Error(err)
	}

	if user != nil {
		return user, nil
	}

	// user not found in cache, try get from db

	exist, err := btx.Db.CheckUserExistByTgId(btx.Id, from.ID)
	if err != nil {
		return nil, err
	}

	if exist {
		user, err = btx.Db.UserByTgId(from.ID, btx.Id)
		if err != nil {
			return nil, err
		}
		// TODO: push user to cache
		return user, nil
	}

	return nil, nil
}

func (btx *TBot) addUser(from *telego.User) (*model.User, error) {
	user := &model.User{
		TgId:      from.ID,
		FirstName: &from.FirstName,
		LastName:  &from.LastName,
		Username:  &from.Username,
		StepId:    1,
		Status:    pgsql.StatusUserActive,
	}

	userID, err := btx.Db.AddUser(btx.Id, user)
	if err != nil {
		return nil, err
	}

	user.Id = userID

	if err = btx.Rdb.SetUser(btx.Id, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (btx *TBot) getComponent(stepID int64) (*model.Component, error) {
	// try get component from cache
	component, err := btx.Rdb.GetComponent(btx.Id, stepID)
	if err != nil {
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

	return nil, nil
}
