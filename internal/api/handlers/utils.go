package handlers

import (
	"github.com/botscubes/bot-service/internal/model"
)

func botComponentCommands(compId int64, commands *[]*model.Command) *[]*command {
	var result = make([]*command, 0)
	for _, c := range *commands {
		if *c.ComponentId == compId {
			cmd := &command{
				Id:         c.Id,
				Type:       c.Type,
				Data:       c.Data,
				NextStepId: c.NextStepId,
			}
			result = append(result, cmd)
		}
	}
	return &result
}

func botComponentsRes(components *[]*model.Component, commands *[]*model.Command) *getBotComponentsRes {
	var result = make(getBotComponentsRes, len(*components))

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
			Commands:   botComponentCommands(v.Id, commands),
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
