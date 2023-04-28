package schema

import "encoding/json"

type NewBotReq struct {
	Title *string `json:"title"`
}

type NewBotRes struct {
	Id int64 `json:"id"`
}

type SetTokenReq struct {
	BotId *json.Number `json:"bot_id"`
	Token *string      `json:"token"`
}

type BotIdReq struct {
	BotId *json.Number `json:"bot_id"`
}
