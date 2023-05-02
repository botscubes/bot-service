package handlers

import (
	"github.com/botscubes/bot-service/internal/model"
)

func botComponentCommands(commands *[]*model.Command) *[]*command {
	var result = make([]*command, len(*commands))
	for i, c := range *commands {
		cmd := &command{
			Id:          c.Id,
			Type:        c.Type,
			Data:        c.Data,
			ComponentId: c.ComponentId,
			NextStepId:  c.NextStepId,
		}
		result[i] = cmd
	}
	return &result
}

func botFullComponentsRes(components *[]*model.ComponentFull) *botFullCompsRes {
	var result = make(botFullCompsRes, len(*components))

	for i, v := range *components {
		cmt := &component{
			Id: v.Id,
			Data: &componentData{
				Type: v.Data.Type,
				Content: &dataContent{
					Text: v.Data.Content.Text,
				},
			},
			Keyboard: &keyboard{
				Buttons: v.Keyboard.Buttons,
			},
			Commands:   botComponentCommands(&v.Commands),
			NextStepId: v.NextStepId,
			Position: &point{
				X: &v.Position.P.X,
				Y: &v.Position.P.Y,
			},
		}

		result[i] = cmt
	}

	return &result
}
