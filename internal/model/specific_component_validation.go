package model

import (
	"errors"
	"strconv"

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

func ValidateSpecificComponentData(ctype string, data map[string]any) *se.ServiceError {
	for key, value := range data {
		validate, ok := SpecificComponentDataValidation[ctype][key]
		if !ok {
			return e.NonExistentParam(key)
		}
		se := validate(value)

		if se != nil {
			return se
		}
	}
	return nil
}

var SpecificComponentDataValidation = map[string]map[string]func(data any) *se.ServiceError{
	"condition": {
		"expression": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("expression")
			}

			return nil
		},
	},
	"message": {
		"text": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("text")
			}
			return nil
		},
	},
	"format": {
		"formatString": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("formatString")
			}
			return nil
		},
	},
	"buttons": {
		"text": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("text")
			}
			return nil
		},
		"buttons": func(data any) *se.ServiceError {
			return nil
		},
	},
	"code": {
		"code": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("code")
			}
			return nil
		},
	},
	"toInt": {
		"source": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("source")
			}
			return nil
		},
		"destination": func(data any) *se.ServiceError {
			_, ok := data.(string)
			if !ok {
				return e.InvalidParam("destination")
			}
			return nil
		},
	},
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
			"idIfError":       true,
			"nextComponentId": true,
			"idIfFalse":       true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
	"message": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError":       true,
			"nextComponentId": true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
	"textInput": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError":       true,
			"nextComponentId": true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
	"format": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError":       true,
			"nextComponentId": true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
	"code": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError":       true,
			"nextComponentId": true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
	"toInt": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError":       true,
			"nextComponentId": true,
		}

		return checkKeyInMap(outputNames, outputName)
	},
	"buttons": func(outputName string) *se.ServiceError {
		outputNames := map[string]bool{
			"idIfError": true,
		}
		err := checkKeyInMap(outputNames, outputName)
		if err == nil {
			return nil
		}

		if _, err := strconv.Atoi(outputName); err != nil {
			return e.OutputPointNameIsNotNumber(outputName)
		}
		return nil
	},
}
