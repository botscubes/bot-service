package errors

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

var (
	ErrInvalidRequest = err.New(1400, "Invalid request")
	ErrInvalidParams  = err.New(1401, "Required parameters are missing")
)

// StartBot handler
var (
	ErrIncorrectTokenFormat = err.New(100, "Token has an incorrect format")
	ErrInvalidToken         = err.New(101, "Invalid token")
	ErrTokenExistInSystem   = err.New(102, "Token already exist in the system")
)

// New new handler
var (
	ErrInvalidTitleLength = err.New(103, "Title is too long")
)
