package handlers

import "encoding/json"

type newBotReq struct {
	UserId *json.Number `json:"user_id"`
	Title  *string      `json:"title"`
}

type newBotRes struct {
	Id int64 `json:"id"`
}

type setTokenReq struct {
	UserId *json.Number `json:"user_id"`
	BotId  *json.Number `json:"bot_id"`
	Token  *string      `json:"token"`
}

type setTokenRes struct {
	Message string `json:"message"`
}
