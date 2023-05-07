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

	// user not found in cache, get from db

	exist, err := btx.Db.CheckUserExistByTgId(btx.Id, from.ID)
	if err != nil {
		return nil, err
	}

	if exist {
		user, err = btx.Db.UserByTgId(from.ID, btx.Id)
		if err != nil {
			return nil, err
		}

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
		Status:    pgsql.StatususerActive,
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
