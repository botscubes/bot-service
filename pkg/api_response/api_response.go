package api_response

import (
	err "github.com/botscubes/user-service/pkg/service_error"
)

type APIResponse struct {
	Ok    bool              `json:"ok"`
	Data  any               `json:"data,omitempty"`
	Error *err.ServiceError `json:"error,omitempty"`
}

func New(ok bool, data any, e *err.ServiceError) *APIResponse {
	return &APIResponse{Ok: ok, Data: data, Error: e}
}
