package model

import (
	e "github.com/botscubes/bot-service/internal/api/errors"
	se "github.com/botscubes/user-service/pkg/service_error"
)

const (
	MaxCommandsCount = 100
)

// Validation command type & data
func (c *CommandParams) Validate() *se.ServiceError {
	if c.Type == nil {
		return e.MissingParam("command.type")
	}

	switch *c.Type {
	case "text":
		return commandTextValidate(c.Data)
	default:
		return e.ErrUnknownCommand
	}
}

// Validation commands list
func (c *CommandsParam) Validate() *se.ServiceError {
	if c == nil {
		return e.MissingParam("commands")
	}

	if len(*c) > MaxCommandsCount {
		return e.ErrTooManyCommands
	}

	for _, v := range *c {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func commandTextValidate(t *string) *se.ServiceError {
	if t == nil {
		return e.MissingParam("command.data")
	}

	if *t == "" {
		return e.InvalidParam("command.data is empty")
	}

	return nil
}

func (r *SetNextStepCommandReq) Validate() *se.ServiceError {
	if r.NextStepId == nil {
		return e.MissingParam("nextStepId")
	}

	if *r.NextStepId < 1 {
		return e.InvalidParam("nextStepId")
	}

	return nil
}
