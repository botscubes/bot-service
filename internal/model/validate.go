package model

import (
	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	se "github.com/botscubes/user-service/pkg/service_error"
)

// Validation of component
func ValidateComponent(d *Data, c *Commands, p *Point) *se.ServiceError {
	// validate data
	if err := d.Validate(); err != nil {
		return err
	}

	// Validate commands
	if err := c.Validate(); err != nil {
		return err
	}

	// Validate position
	return p.Validate()
}

// Validation of component type and data value by type
func (d *Data) Validate() *se.ServiceError {
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

// Validation of component commands list
func (c *Commands) Validate() *se.ServiceError {
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

// Validation of component command
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

// Validation of component position
func (p *Point) Validate() *se.ServiceError {
	if p == nil {
		return e.InvalidParam("position")
	}

	if int64(p.X) < 0 || int64(p.X) > config.MaxPositionX {
		return e.IncorrectVal("position.x")
	}

	if int64(p.Y) < 0 || int64(p.Y) > config.MaxPositionY {
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
