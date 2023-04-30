package model

import (
	"github.com/jackc/pgx/pgtype"
)

type Component struct {
	Id       int64         `json:"id"`
	Data     *Data         `json:"data"`
	Keyboard *Keyboard     `json:"keyboard,omitempty"`
	NextId   *int64        `json:"nextId,omitempty"`
	Position *pgtype.Point `json:"position"`
	Status   int           `json:"status"`
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
	Id     *int64  `json:"id"`
	Type   *string `json:"type"`
	Data   *string `json:"data"`
	NextId *int64  `json:"nextId,omitempty"`
}
