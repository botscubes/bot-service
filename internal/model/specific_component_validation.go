package model

import (
	"errors"
)

func validateId(data any) error {
	v, ok := data.(int64)
	if !ok {
		return errors.New("The value must be an integer value")
	}
	if v < 0 {
		return errors.New("The number must not be negative")
	}

	return nil
}

var specific_component_validation = map[string]map[string]func(data any) error{
	"start": {
		"nextComponentId": validateId,
	},
}
