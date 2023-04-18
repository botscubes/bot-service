package errors

// TOOD: create Responce package with struct
//  Ok bool
//  Data
//  Error *service_error

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

var (
	ErrInvalidRequest = err.New(1400, "Invalid request")
)

var (
	ErrIncorrectTokenFormat = err.New(100, "Token has an incorrect format")
	ErrInvalidToken         = err.New(101, "Invalid token")
	ErrTokenExistInSystem   = err.New(102, "Token already exist in the system")
)

var (
	Success = err.New(0, "Success")
)
