package handlers

import (
	"github.com/botscubes/bot-service/internal/model"
)

func componentCommands(commands *[]*model.Command) *[]*command {
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

func dataContentsRes(contents *[]*model.CompContent) *[]*dataContent {
	var result = make([]*dataContent, len(*contents))

	for i, c := range *contents {
		cnt := &dataContent{
			Text: c.Text,
		}
		result[i] = cnt
	}

	return &result
}

func componentRes(v *model.Component) *component {
	return &component{
		Id: v.Id,
		Data: &componentData{
			Type:    v.Data.Type,
			Content: dataContentsRes(&v.Data.Content),
		},
		Keyboard: &keyboard{
			Buttons: v.Keyboard.Buttons,
		},
		Commands:   componentCommands(&v.Commands),
		NextStepId: v.NextStepId,
		IsMain:     v.IsMain,
		Position: &point{
			X: &v.Position.P.X,
			Y: &v.Position.P.Y,
		},
	}
}

func componentsRes(components *[]*model.Component) *getBotCompsRes {
	var result = make(getBotCompsRes, len(*components))

	for i, v := range *components {
		result[i] = componentRes(v)
	}

	return &result
}

func dataContentsMod(contents *[]*dataContent) *[]*model.CompContent {
	var result = make([]*model.CompContent, len(*contents))

	for i, c := range *contents {
		cnt := &model.CompContent{
			Text: c.Text,
		}
		result[i] = cnt
	}

	return &result
}
