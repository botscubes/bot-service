package model

import se "github.com/botscubes/user-service/pkg/service_error"
import e "github.com/botscubes/bot-service/internal/api/errors"

func (c *Connection) Validate() *se.ServiceError {
	if c.SourceComponentId == nil {
		return e.MissingParam("sourceComponentId")
	}
	if c.SourcePointName == nil {
		return e.MissingParam("sourcePointName")
	}
	if c.TargetComponentId == nil {
		return e.MissingParam("targetComponentId")
	}
	if c.RelativePointPosition == nil {
		return e.MissingParam("relativePointPosition")
	}

	if *c.SourceComponentId < 0 {
		return e.InvalidParam("sourceComponentId must not be negative")
	}

	if *c.TargetComponentId < 0 {
		return e.InvalidParam("targetComponentId must not be negative")
	}

	return c.RelativePointPosition.Validate()
}
