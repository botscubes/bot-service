package model

import se "github.com/botscubes/user-service/pkg/service_error"
import e "github.com/botscubes/bot-service/internal/api/errors"

func (c *SourceConnectionPoint) Validate(componentType string) *se.ServiceError {
	if c.SourceComponentId == nil {
		return e.MissingParam("sourceComponentId")
	}
	if c.SourcePointName == nil {
		return e.MissingParam("sourcePointName")
	}
	if *c.SourceComponentId < 0 {
		return e.InvalidParam("sourceComponentId must not be negative")
	}
	validate, ok := SpecificComponentOutputValidation[componentType]
	if !ok {
		return e.ErrValidation
	}
	se := validate(*c.SourcePointName)
	if se != nil {
		return se
	}
	return nil
}

func (c *Connection) Validate(componentType string) *se.ServiceError {

	if c.TargetComponentId == nil {
		return e.MissingParam("targetComponentId")
	}
	if c.RelativePointPosition == nil {
		return e.MissingParam("relativePointPosition")
	}

	if *c.TargetComponentId < 0 {
		return e.InvalidParam("targetComponentId must not be negative")
	}
	if err := c.RelativePointPosition.Validate(); err != nil {
		return err
	}
	return c.SourceConnectionPoint.Validate(componentType)
}
