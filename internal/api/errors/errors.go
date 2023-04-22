package errors

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

var (
	ErrInvalidRequest = err.New(1400, "Invalid request")
	ErrInvalidParams  = err.New(1401, "Required parameters are missing")
)

var (
	ErrIncorrectTokenFormat  = err.New(100, "Token has an incorrect format")
	ErrInvalidToken          = err.New(101, "Invalid token")
	ErrTokenExistInSystem    = err.New(102, "Token already exist in the system")
	ErrInvalidTitleLength    = err.New(103, "Title is too long")
	ErrBotNotFound           = err.New(104, "Bot not found")
	ErrTokenAlreadyInstalled = err.New(105, "Token is already installed")
	ErrTokenAlreadyExists    = err.New(106, "Token already exists")
)
