package bot

import (
	"errors"

	"github.com/botscubes/bot-service/internal/config"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/model"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (btx *TBot) mainHandler() th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		stepID, err := btx.getUserStep(update.Message.From)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				log.Error(err)
				return
			}

			if err = btx.addUser(update.Message.From); err != nil {
				log.Error(err)
				return
			}

			stepID = config.MainComponentId
		}

		// find next component for execute

		var origComponent *model.Component
		var component *model.Component
		origStepID := stepID

		// for cycle detect
		stepsPassed := make(map[int64]struct{})
		isFound := false

		for {
			// This part of the loop (before the "isFound" condition) is used to automatically
			// skip the starting component and undefined components.
			// Also, the next component is selected here by the id found in the second part of the cycle.

			if _, ok := stepsPassed[stepID]; ok {
				log.Warnf("cycle detected: bot #%d", btx.Id)
				return
			}

			stepsPassed[stepID] = struct{}{}

			component, err = btx.getComponent(stepID)
			if err != nil {
				if errors.Is(err, rdb.ErrNotFound) {
					stepID = config.MainComponentId
					continue
				}

				log.Error(err)
				return
			}

			if origComponent == nil {
				origComponent = component
			}

			// check main component
			if component.IsMain {
				if component.NextStepId == nil {
					log.Warnf("main component does not have link to the following: bot #%q", btx.Id)
					return
				}

				stepID = *component.NextStepId
				continue
			}

			if isFound {
				// next component was found successfully
				break
			}

			// In this part, the id of the next component is determined.
			// In case of successful identification of the ID, an additional check occurs in the first part of the cycle.

			isFound = true

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
			stepID = origStepID
			break
		}

		if err := btx.Rdb.SetUserStep(btx.Id, update.Message.From.ID, stepID); err != nil {
			log.Error(err)
		}

		if err := exec(bot, &update, component.Data); err != nil {
			log.Error(err)
		}
	}
}

func exec(bot *telego.Bot, update *telego.Update, data *model.Data) error {
	switch *data.Type {
	case "text":
		_, err := bot.SendMessage(tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			*(*data.Content)[0].Text,
		))
		if err != nil {
			return err
		}
	default:
		log.Warn("Unknown type method: ", *data.Type)
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
