package errors

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

var (
	ErrInvalidRequest = err.New(1400, "Invalid request")
	ErrInvalidParams  = err.New(1411, "Required parameters are missing")
	ErrInternalServer = err.New(1500, "Internal server error")
	ErrUnauthorized   = err.New(1401, "Unauthorized")
)

var (
	ErrIncorrectTokenFormat  = err.New(100, "Token has an incorrect format")
	ErrInvalidToken          = err.New(101, "Invalid token")
	ErrInvalidTitleLength    = err.New(103, "Title is too long")
	ErrBotNotFound           = err.New(104, "Bot not found")
	ErrTokenAlreadyInstalled = err.New(105, "The bot already has a token")
	ErrTokenAlreadyExists    = err.New(106, "Token already exists in system")
	ErrTokenNotFound         = err.New(107, "Token not found")
	ErrStartBot              = err.New(108, "Start bot error")
	ErrStopBot               = err.New(109, "Stop bot error")
	ErrBotNotRunning         = err.New(110, "The bot is not running")
)
