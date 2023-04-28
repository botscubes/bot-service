package handlers

import "encoding/json"

type newBotReq struct {
	Title *string `json:"title"`
}

type newBotRes struct {
	Id int64 `json:"id"`
}

type setTokenReq struct {
	BotId *json.Number `json:"bot_id"`
	Token *string      `json:"token"`
}

type botIdReq struct {
	BotId *json.Number `json:"bot_id"`
}
