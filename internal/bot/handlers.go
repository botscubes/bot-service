package bot

import (
	ct "github.com/botscubes/bot-service/internal/components"
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
			btx.Users[chatID] = &ct.User{
				StepId: 1,
			}
		}

		stepID := &btx.Users[chatID].StepId

		log.Info(*stepID)

		// check start component
		if btx.Components[*stepID].IsMain {
			btx.Users[chatID].StepId = *btx.Components[*stepID].NextStepId
		}

		command := determineCommand(&update.Message.Text, btx.Components[*stepID].Commands)
		if command != nil {
			btx.Users[chatID].StepId = *command.NextStepId
		}

		log.Info(*stepID)

		execMethod(bot, &update, btx.Components[*stepID].Data)
	}
}

func execMethod(bot *telego.Bot, update *telego.Update, data *ct.Data) {
	switch *data.Type {
	case "text":
		_, _ = bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			*(*data.Content)[0].Text,
		))
	}
}

func determineCommand(mes *string, commands *[]*ct.Command) *ct.Command {
	// work with command type - text

	for _, command := range *commands {
		if *command.Data == *mes {
			return command
		}
	}

	return nil
}
