package model

import (
	"github.com/jackc/pgx/pgtype"
)

type Component struct {
	Id       int64
	Data     Data
	Keyboard Keyboard
	NextId   *int64
	Position pgtype.Point
	Status   int
}

type Data struct {
	Type    string  `json:"type,omitempty"`
	Content Content `json:"content,omitempty"`
}

type Content struct {
	Text string `json:"text,omitempty"`
}

type Keyboard struct {
	Buttons [][]Button `json:"buttons,omitempty"`
}

type Button struct {
	Text   string `json:"text,omitempty"`
	NextId int64  `json:"next_id,omitempty"`
}
