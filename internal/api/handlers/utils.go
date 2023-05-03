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

func compContentsRes(contents *[]*model.CompContent) *[]*dataContent {
	var result = make([]*dataContent, len(*contents))

	for i, c := range *contents {
		cnt := &dataContent{
			Text: c.Text,
		}
		result[i] = cnt
	}

	return &result
}

func botComponentRes(v *model.Component) *component {

	cmt := &component{
		Id: v.Id,
		Data: &componentData{
			Type:    v.Data.Type,
			Content: compContentsRes(&v.Data.Content),
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

func compContentsMod(contents *[]*dataContent) *[]*model.CompContent {
	var result = make([]*model.CompContent, len(*contents))

	for i, c := range *contents {
		cnt := &model.CompContent{
			Text: c.Text,
		}
		result[i] = cnt
	}

	return &result
}
