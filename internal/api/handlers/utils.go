package handlers

// TODO: create pkg

import (
	ct "github.com/botscubes/bot-service/internal/components"
	"github.com/botscubes/bot-service/internal/model"
)

func componentCommands(commands *[]*model.Command) *[]*ct.Command {
	var result = make([]*ct.Command, len(*commands))

	for i, c := range *commands {
		cmd := &ct.Command{
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

func dataContentsRes(contents *[]*model.CompContent) *[]*ct.Content {
	var result = make([]*ct.Content, len(*contents))

	for i, c := range *contents {
		cnt := &ct.Content{
			Text: c.Text,
		}
		result[i] = cnt
	}

	return &result
}

func componentRes(v *model.Component) *ct.Component {
	return &ct.Component{
		Id: &v.Id,
		Data: &ct.Data{
			Type:    v.Data.Type,
			Content: dataContentsRes(&v.Data.Content),
		},
		Keyboard: &ct.Keyboard{
			Buttons: v.Keyboard.Buttons,
		},
		Commands:   componentCommands(&v.Commands),
		NextStepId: v.NextStepId,
		IsMain:     v.IsMain,
		Position: &ct.Point{
			X: &v.Position.P.X,
			Y: &v.Position.P.Y,
		},
	}
}

// TODO: Add converts diffrent funcs for Editor and Bot

func componentsRes(components *[]*model.Component) *[]*ct.Component {
	var result = make([]*ct.Component, len(*components))

	for i, v := range *components {
		result[i] = componentRes(v)
	}

	return &result
}

func dataContentsMod(contents *[]*ct.Content) *[]*model.CompContent {
	var result = make([]*model.CompContent, len(*contents))

	for i, c := range *contents {
		cnt := &model.CompContent{
			Text: c.Text,
		}
		result[i] = cnt
	}

	return &result
}
