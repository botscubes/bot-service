package model

import (
	"errors"

	e "github.com/botscubes/bot-service/internal/api/errors"

	se "github.com/botscubes/user-service/pkg/service_error"
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

var SpecificComponentDataValidation = map[string]map[string]func(data any) error{
	"start": {},
}

func checkKeyInMap(m map[string]bool, k string) *se.ServiceError {
	_, ok := m[k]
	if !ok {
		return e.NoOutputPointName(k)
	}
	return nil
}

var SpecificComponentPointValidation = map[string]func(pointName string) *se.ServiceError{
	"start": func(pointName string) *se.ServiceError {
		pointNames := map[string]bool{
			"nextComponentId": true,
		}
		return checkKeyInMap(pointNames, pointName)
	},
	"condition": func(pointName string) *se.ServiceError {
		pointNames := map[string]bool{
			"idIfError": true,
		}

		return checkKeyInMap(pointNames, pointName)
	},
}
