package model

// May be create package "components" -_-.._-

import (
	"unicode/utf8"

	e "github.com/botscubes/bot-service/internal/api/errors"
	se "github.com/botscubes/user-service/pkg/service_error"
)

const (
	maxMessageLen = 4096
)

func startComponentValidate() *se.ServiceError {
	return e.ErrMainComponent
}

func textComponentValidate(c *[]*Content) *se.ServiceError {
	if len(*c) != 1 {
		return e.InvalidParam("data.content")
	}

	if (*c)[0].Text == nil {
		return e.MissingParam("data.content.text")
	}

	// check text min length
	if utf8.RuneCountInString(*(*c)[0].Text) < 1 {
		return e.ErrComponentTextTooShort
	}

	// check text max length
	if utf8.RuneCountInString(*(*c)[0].Text) > maxMessageLen {
		return e.ErrComponentTextTooLong
	}

	return nil
}
