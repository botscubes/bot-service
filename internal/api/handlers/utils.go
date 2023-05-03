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

func botComponentRes(v *model.Component) *component {
	var dcontent dataContent

	if v.Data.Content != nil {
		dcontent.Text = v.Data.Content.Text
	}

	cmt := &component{
		Id: v.Id,
		Data: &componentData{
			Type:    v.Data.Type,
			Content: &dcontent,
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

	return cmt
}

func botComponentsRes(components *[]*model.Component) *getBotCompsRes {
	var result = make(getBotCompsRes, len(*components))

	for i, v := range *components {
		result[i] = botComponentRes(v)
	}

	return &result
}
