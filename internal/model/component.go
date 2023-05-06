package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Component struct {
	Id         int64         `json:"id"`
	Data       *Data         `json:"data"`
	Keyboard   *Keyboard     `json:"keyboard,omitempty"`
	Commands   []*Command    `json:"commands,omitempty"`
	NextStepId *int64        `json:"next_step_id,omitempty"`
	IsMain     bool          `json:"is_main"`
	Position   *pgtype.Point `json:"position"`
	Status     int           `json:"status"`
}

type Data struct {
	Type    *string        `json:"type"`
	Content []*CompContent `json:"content"`
}

type CompContent struct {
	Text *string `json:"text,omitempty"`
}

type Keyboard struct {
	Buttons [][]*int64 `json:"buttons"`
}

type Command struct {
	Id          *int64  `json:"id"`
	Type        *string `json:"type"`
	Data        *string `json:"data"`
	ComponentId *int64  `json:"component_id"`
	NextStepId  *int64  `json:"next_step_id"`
	Status      int     `json:"status"`
}
