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

var SpecificComponentOutputValidation = map[string]func(outputName string) *se.ServiceError{
	"start": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"nextComponentId": true,
		}
		return checkKeyInMap(outputNames, outputName)
	},
	"condition": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError": true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
}
