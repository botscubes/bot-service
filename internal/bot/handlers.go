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
		user, err := btx.getUser(update.Message.From)
		if err != nil {
			log.Error(err)
			return
		}

		if user == nil {
			user, err = btx.addUser(update.Message.From)
			if err != nil {
				log.Error(err)
				return
			}
		}

		stepID := user.StepId

		log.Info(btx.Id, stepID)
		btx.Rdb.PrintAllComponents(btx.Id)
		// move to separate func
		component, err := btx.Rdb.GetComponent(btx.Id, stepID)
		if err != nil {
			log.Error(err)
			// try get from db
		}

		if component == nil {
			log.Debug("Component not found")
			return
		}

		// todo: check nextStepid = nil

		log.Debug("ASd")
		log.Debugf("%+v", *component)
		// check start component
		if component.IsMain {
			stepID = *component.NextStepId
		}

		command := determineCommand(&update.Message.Text, component.Commands)
		if command != nil {
			stepID = *command.NextStepId
			user.StepId = stepID

			// TODO: save in cache only stepId
		}

		log.Info(stepID)
		user.StepId = stepID
		btx.Rdb.SetUser(btx.Id, user)

		// TODO: get component

		if err := execMethod(bot, &update, component.Data); err != nil {
			log.Error(err)
		}
	}
}

func execMethod(bot *telego.Bot, update *telego.Update, data *model.Data) error {
	switch *data.Type {
	case "text":
		e, err := bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			*(*data.Content)[0].Text,
		))
		if err != nil {
			return err
		}

		log.Warn(e)
	}

	return nil
}

// Determine commnad by !message text!
func determineCommand(mes *string, commands *model.Commands) *model.Command {
	// work with command type - text

	for _, command := range *commands {
		if *command.Data == *mes {
			return command
		}
	}

	return nil
}
