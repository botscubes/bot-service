package model

import (
	"github.com/botscubes/bot-components/components"
	e "github.com/botscubes/bot-service/internal/api/errors"
	"github.com/botscubes/bot-service/internal/config"
	se "github.com/botscubes/user-service/pkg/service_error"
)

const (
	MaxPositionX = 10000
	MaxPositionY = 10000
)

// Validation of component data
func (d *ComponentData) Validate() *se.ServiceError {
	if d == nil {
		return e.MissingParam("data")
	}

	if d.Type == nil {
		return e.MissingParam("data.type")
	}

	if d.Content == nil {
		return e.MissingParam("data.content")
	}

	switch *d.Type {
	case "start":
		return startComponentValidate()
	case "text":
		return textComponentValidate(d.Content)
	default:
		return e.ErrUnknownComponent
	}
}

// Validation of component position
func (p *Point) Validate() *se.ServiceError {
	if p == nil {
		return e.MissingParam("position")
	}

	//if int64(p.X) < 0 || int64(p.X) > MaxPositionX {
	//		return e.InvalidParam("position.x")
	//	}

	//	if int64(p.Y) < 0 || int64(p.Y) > MaxPositionY {
	//		return e.InvalidParam("position.y")
	//	}

	return nil
}

func (r *AddComponentReq) Validate() *se.ServiceError {
	if r.Type != components.TypeFormat &&
		r.Type != components.TypeCondition &&
		r.Type != components.TypeMessage &&
		r.Type != components.TypeTextInput &&
		r.Type != components.TypeCode &&
		r.Type != components.TypeButtons &&
		r.Type != components.TypeToInt &&
		r.Type != components.TypeMove &&
		r.Type != components.TypeHTTP &&
		r.Type != components.TypeFromJSON &&
		r.Type != components.TypePhoto {
		return e.InvalidParam("type: " + r.Type)
	}

	// Validate position
	return r.Position.Validate()
}

func (r *SetNextStepComponentReq) Validate() *se.ServiceError {
	if r.NextStepId == nil {
		return e.MissingParam("nextStepId")
	}

	if *r.NextStepId < 1 {
		return e.InvalidParam("nextStepId")
	}

	return nil
}

func (r *DelSetComponentsReq) Validate() *se.ServiceError {
	if r.Data == nil {
		return e.MissingParam("data")
	}

	if len(*r.Data) == 0 {
		return e.InvalidParam("data")
	}

	// check main component in req
	for _, v := range *r.Data {
		if v == config.MainComponentId {
			return e.ErrMainComponent
		}
	}

	return nil
}

func (r *UpdComponentReq) Validate() *se.ServiceError {
	if r.Data == nil && r.Position == nil {
		return e.ErrBadRequest
	}

	if r.Data != nil {
		if err := r.Data.Validate(); err != nil {
			return err
		}
	}

	if r.Position != nil {
		if err := r.Position.Validate(); err != nil {
			return err
		}
	}

	return nil
}
