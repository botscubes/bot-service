package bot

import (
	bu "github.com/botscubes/bot-service/internal/bot/util"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendMessage(bot *telego.Bot, update *telego.Update, component *model.Component) error {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		*(*component.Data.Content)[0].Text,
	)

	if len(*component.Commands) > 0 {
		message.WithReplyMarkup(bu.Keyboard(component.Commands, component.Keyboard).WithResizeKeyboard())
	}

	_, err := bot.SendMessage(message)
	return err
}
