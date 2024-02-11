package errors

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

var (
	BadRequestCode     = 1400
	NotFoundCode       = 1404
	InvalidParamCode   = 1411
	MissingParamCode   = 1412
	InternalServerCode = 1500
	UnauthorizedCode   = 1401
	IncorrectValCode   = 1402
)

var (
	ErrBadRequest     = err.New(BadRequestCode, "Bad request")
	ErrNotFound       = err.New(NotFoundCode, "Not Found")
	ErrInvalidParam   = err.New(InvalidParamCode, "Invalid parameter value")
	ErrMissingParam   = err.New(MissingParamCode, "Required parameter is missing")
	ErrInternalServer = err.New(InternalServerCode, "Internal server error")
	ErrUnauthorized   = err.New(UnauthorizedCode, "Unauthorized")
	ErrIncorrectVal   = err.New(IncorrectValCode, "Incorrect value")
)

var (
	ErrIncorrectTokenFormat  = err.New(100, "Token has an incorrect format")
	ErrInvalidToken          = err.New(101, "Invalid token")
	ErrTitleTooLong          = err.New(102, "Title is too long")
	ErrBotNotFound           = err.New(103, "Bot not found")
	ErrTokenAlreadyInstalled = err.New(104, "The bot already has a token")
	ErrTokenAlreadyExists    = err.New(105, "Token already exists in system")
	ErrTokenNotFound         = err.New(106, "Token not found")
	ErrStartBot              = err.New(107, "Start bot error")
	ErrStopBot               = err.New(108, "Stop bot error")
	ErrBotAlreadyStopped     = err.New(109, "The bot already stopped")
	ErrComponentNotFound     = err.New(110, "Component not found")
	ErrNextComponentNotFound = err.New(111, "Next component not found")
	ErrCommandNotFound       = err.New(112, "Command not found")
	ErrMainComponent         = err.New(113, "The action is not available for the main component")
	ErrUnknownComponent      = err.New(114, "Unknown component")
	ErrUnknownCommand        = err.New(115, "Unknown command")
	ErrBotAlreadyRunning     = err.New(116, "The bot already running")
	ErrBotNeedsStopped       = err.New(117, "Bot needs to be stopped")
	ErrNewBot                = err.New(118, "Create bot error")
	ErrTitleTooShort         = err.New(119, "Title is too short")
	ErrComponentTextTooShort = err.New(120, "Text is too short")
	ErrComponentTextTooLong  = err.New(121, "Text is too long")
	ErrTooManyCommands       = err.New(122, "Too many commands")
	ErrGroupNotFound         = err.New(123, "Group not found")
	ErrDeleteStartComponent  = err.New(124, "Starting component cannot be deleted")
)

func InvalidParam(mes string) *err.ServiceError {
	return err.New(InvalidParamCode, "Invalid parameter value: "+mes)
}

func MissingParam(mes string) *err.ServiceError {
	return err.New(InvalidParamCode, "Required parameter is missing: "+mes)
}

func NoOutputPointName(mes string) *err.ServiceError {
	return err.New(InvalidParamCode, "No output point name: "+mes)
}
