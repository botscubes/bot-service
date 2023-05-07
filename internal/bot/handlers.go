package bot

import (
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: may be -> self-contained functions (without (btx *TBot) )

func (btx *TBot) mainHandler() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)

		if _, ok := btx.Users[chatID]; !ok {
			btx.Users[chatID] = &model.User{
				StepId: 1,
			}
		}

		stepID := &btx.Users[chatID].StepId

		log.Info(*stepID)

		component, err := btx.Rdb.GetComponet(btx.Id, *stepID)
		if err != nil {
			log.Error(err)
		}

		// check start component
		if component.IsMain {
			btx.Users[chatID].StepId = *component.NextStepId
		}

		command := determineCommand(&update.Message.Text, component.Commands)
		if command != nil {
			btx.Users[chatID].StepId = *command.NextStepId
		}

		log.Info(*stepID)

		execMethod(bot, &update, component.Data)
	}
}

func execMethod(bot *telego.Bot, update *telego.Update, data *model.Data) {
	switch *data.Type {
	case "text":
		_, _ = bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			*(*data.Content)[0].Text,
		))
	}
}

func determineCommand(mes *string, commands *model.Commands) *model.Command {
	// work with command type - text

	for _, command := range *commands {
		if *command.Data == *mes {
			return command
		}
	}

	return nil
}
