package errors

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

var (
	InvalidRequestCode = 1400
	InvalidParamCode   = 1411
	InternalServerCode = 1500
	UnauthorizedCode   = 1401
	IncorrectValCode   = 1402
)

var (
	ErrInvalidRequest = err.New(InvalidRequestCode, "Invalid request")
	ErrInvalidParam   = err.New(InvalidParamCode, "Required parameter is missing")
	ErrInternalServer = err.New(InternalServerCode, "Internal server error")
	ErrUnauthorized   = err.New(UnauthorizedCode, "Unauthorized")
	ErrIncorrectVal   = err.New(IncorrectValCode, "Incorrect value")
)

var (
	ErrIncorrectTokenFormat  = err.New(100, "Token has an incorrect format")
	ErrInvalidToken          = err.New(101, "Invalid token")
	ErrInvalidTitleLength    = err.New(102, "Title is too long")
	ErrBotNotFound           = err.New(103, "Bot not found")
	ErrTokenAlreadyInstalled = err.New(104, "The bot already has a token")
	ErrTokenAlreadyExists    = err.New(105, "Token already exists in system")
	ErrTokenNotFound         = err.New(106, "Token not found")
	ErrStartBot              = err.New(107, "Start bot error")
	ErrStopBot               = err.New(108, "Stop bot error")
	ErrBotNotRunning         = err.New(109, "The bot is not running")
	ErrComponentNotFound     = err.New(110, "Component not found")
	ErrNextComponentNotFound = err.New(111, "Next component not found")
	ErrCommandNotFound       = err.New(112, "Command not found")
	ErrMainComponent         = err.New(113, "The action is not available for the main component")
	ErrUnknownComponent      = err.New(114, "Unknown component")
	ErrUnknownCommand        = err.New(115, "Unknown command")
)

func InvalidParam(mes string) *err.ServiceError {
	return err.New(InvalidParamCode, "Required parameter is missing: "+mes)
}

func IncorrectVal(mes string) *err.ServiceError {
	return err.New(IncorrectValCode, "Incorrect value: "+mes)
}
