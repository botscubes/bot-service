package model

import (
	"regexp"
	"unicode/utf8"

	e "github.com/botscubes/bot-service/internal/api/errors"
	se "github.com/botscubes/user-service/pkg/service_error"
)

const (
	MaxTitleLen = 50                     // Max bot title length
	tokenRegexp = `^\d{9,10}:[\w-]{35}$` //nolint:gosec
)

func (r *NewBotReq) Validate() *se.ServiceError {
	if r.Title == nil || *r.Title == "" {
		return e.MissingParam("title")
	}

	// check title min length
	if utf8.RuneCountInString(*r.Title) < 1 {
		return e.ErrTitleTooShort
	}

	// check title max length
	if utf8.RuneCountInString(*r.Title) > MaxTitleLen {
		return e.ErrTitleTooLong
	}

	return nil
}

func (r *SetBotTokenReq) Validate() *se.ServiceError {

	if r.Token == nil {
		return e.MissingParam("token")
	}

	reg := regexp.MustCompile(tokenRegexp)
	if !reg.MatchString(*r.Token) {
		return e.ErrIncorrectTokenFormat
	}

	return nil
}
