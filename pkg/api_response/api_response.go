package api_response

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

type ApiResponse struct {
	Ok    bool              `json:"ok"`
	Data  interface{}       `json:"data,omitempty"`
	Error *err.ServiceError `json:"error,omitempty"`
}

func New(ok bool, data interface{}, e *err.ServiceError) *ApiResponse {
	return &ApiResponse{Ok: ok, Data: data, Error: e}
}
