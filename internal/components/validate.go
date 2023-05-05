package components

import (
	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	se "github.com/botscubes/user-service/pkg/service_error"
)

func CheckIsMain(id int64) bool {
	return id == config.MainComponentId
}

func ValidateComponent(d *Data, c *[]*Command, p *Point) *se.ServiceError {
	if err := ValidateData(d); err != nil {
		return err
	}

	if err := ValidateCommands(c); err != nil {
		return err
	}

	return ValidatePosition(p)
}

func ValidateData(d *Data) *se.ServiceError {
	if d == nil {
		return e.InvalidParam("data")
	}

	if d.Type == nil {
		return e.InvalidParam("data.type")
	}

	if d.Content == nil {
		return e.InvalidParam("data.content")
	}

	switch *d.Type {
	case "start":
		return vStart()
	case "text":
		return vContentText(d.Content)
	default:
		return e.ErrUnknownComponent
	}
}

func ValidateCommands(c *[]*Command) *se.ServiceError {
	if c == nil {
		return e.InvalidParam("commands")
	}

	for _, v := range *c {
		if err := ValidateCommand(v.Type, v.Data); err != nil {
			return err
		}
	}

	return nil
}

func ValidateCommand(t, d *string) *se.ServiceError {
	if t == nil {
		return e.InvalidParam("command.type")
	}

	switch *t {
	case "text":
		return vCommandText(d)
	default:
		return e.ErrUnknownCommand
	}
}

func ValidatePosition(p *Point) *se.ServiceError {
	if p == nil {
		return e.InvalidParam("position")
	}

	if p.X == nil {
		return e.InvalidParam("position.x")
	}

	if p.Y == nil {
		return e.InvalidParam("position.y")
	}

	if int64(*p.X) < 0 || int64(*p.X) > config.MaxPositionX {
		return e.IncorrectVal("position.x")
	}

	if int64(*p.Y) < 0 || int64(*p.Y) > config.MaxPositionY {
		return e.IncorrectVal("position.y")
	}

	return nil
}

func vStart() *se.ServiceError {
	return e.ErrMainComponent
}

func vContentText(c *[]*Content) *se.ServiceError {
	if len(*c) != 1 {
		return e.IncorrectVal("data.content")
	}

	if (*c)[0].Text == nil {
		return e.InvalidParam("data.content.text")
	}

	if *(*c)[0].Text == "" {
		return e.IncorrectVal("data.content.text is empty")
	}

	return nil
}

func vCommandText(t *string) *se.ServiceError {
	if t == nil {
		return e.InvalidParam("command.data")
	}

	if *t == "" {
		return e.IncorrectVal("command.data is empty")
	}

	return nil
}
