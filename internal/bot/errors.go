package bot

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrBotNotFound = errors.New("bot not found")
	ErrTgAuth401   = errors.New(`telego: health check: telego: getMe(): api: 401 "Unauthorized"`)
)
