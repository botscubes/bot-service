package model

import (
	"github.com/jackc/pgx/pgtype"
)

type Component struct {
	Id         int64         `json:"id"`
	Data       *Data         `json:"data"`
	Keyboard   *Keyboard     `json:"keyboard,omitempty"`
	NextStepId *int64        `json:"nextStepId,omitempty"`
	Position   *pgtype.Point `json:"position"`
	Status     int           `json:"status"`
}

type Data struct {
	Type    *string  `json:"type"`
	Content *Content `json:"content"`
}

type Content struct {
	Text *string `json:"text,omitempty"`
}

type Keyboard struct {
	Buttons [][]*int64 `json:"buttons"`
}

type Command struct {
	Id          *int64  `json:"id"`
	Type        *string `json:"type"`
	Data        *string `json:"data"`
	ComponentId *int64  `json:"componentId"`
	NextStepId  *int64  `json:"nextStepId,omitempty"`
	Status      int     `json:"status"`
}
