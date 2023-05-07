package bot

import (
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (btx *TBot) mainHandler() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		user, err := btx.getUser(update.Message.From)
		if err != nil {
			log.Error(err)
			return
		}

		// user not found
		if user == nil {
			user, err = btx.addUser(update.Message.From)
			if err != nil {
				log.Error(err)
				return
			}
		}

		// find next component for execute

		var origComponent *model.Component
		var component *model.Component

		// for cycle detect
		stepsPassed := make(map[int64]struct{})
		stepID := user.StepId
		firstCheck := true
		for {
			// This part of the loop (before the "firstCheck" condition) is used to automatically
			// skip the starting component and undefined components.
			// Also, the next component is selected here by the id found in the second part of the cycle.

			if _, ok := stepsPassed[stepID]; ok {
				log.Warnf("cycle detected: bot #%d", btx.Id)
				return
			}

			stepsPassed[stepID] = struct{}{}

			component, err = btx.getComponent(stepID)
			if err != nil {
				log.Error(err)
				return
			}

			//  if component is nil, run start component
			if component == nil {
				stepID = 1
				// user.StepId = stepID
				// btx.Rdb.SetUser(btx.Id, user)
				continue
			}

			if origComponent == nil {
				origComponent = component
			}

			// check start component
			if component.IsMain {
				if component.NextStepId == nil {
					log.Warnf("start component does not have link to the following: bot #%q", btx.Id)
					return
				}

				stepID = *component.NextStepId
				// user.StepId = stepID
				// btx.Rdb.SetUser(btx.Id, user)
				continue
			}

			if !firstCheck {
				// next component was found successfully
				break
			}

			// In this part, the id of the next component is determined.
			// In case of successful identification of the ID, an additional check occurs in the first part of the cycle.

			firstCheck = false

			if component.NextStepId != nil {
				stepID = *component.NextStepId
				continue
			}

			command := determineCommand(&update.Message.Text, component.Commands)
			if command != nil && command.NextStepId != nil {
				stepID = *command.NextStepId
				// user.StepId = stepID
				continue
			}

			// next component not found, will be executed initial (current) component
			component = origComponent
			break
		}

		user.StepId = 1
		if err := btx.Rdb.SetUser(btx.Id, user); err != nil {
			log.Error(err)
		}

		if err := execMethod(bot, &update, component.Data); err != nil {
			log.Error(err)
		}
	}
}

func execMethod(bot *telego.Bot, update *telego.Update, data *model.Data) error {

	// for linter (switch with one case)
	if *data.Type == "text" {
		e, err := bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			*(*data.Content)[0].Text,
		))
		if err != nil {
			return err
		}

		log.Warn(e)
	}

	// switch *data.Type {
	// case "text":
	// 	e, err := bot.SendMessage(tu.Messagef(
	// 		tu.ID(update.Message.Chat.ID),
	// 		*(*data.Content)[0].Text,
	// 	))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	log.Warn(e)
	// }

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
