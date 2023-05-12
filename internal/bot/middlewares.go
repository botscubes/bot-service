package bot

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func (btx *TBot) regUserMW(next th.Handler) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		var user *telego.User

		// Get user ID from the message
		if update.Message != nil && update.Message.From != nil {
			user = update.Message.From
		}

		// Get user ID from the callback query
		if update.CallbackQuery != nil {
			user = &update.CallbackQuery.From
		}

		// check user exist in cache
		ex, err := btx.Rdb.CheckUserExist(btx.Id, user.ID)
		if err != nil {
			btx.log.Error(err)
		}

		// user not found in cache, check db
		if ex == 0 {
			exist, err := btx.Db.CheckUserExistByTgId(btx.Id, user.ID)
			if err != nil {
				btx.log.Error(err)
				return
			}

			if !exist {
				if err = btx.addUser(user); err != nil {
					return
				}
			}
		}

		next(bot, update)
	}
}
