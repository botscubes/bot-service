package schema

import "encoding/json"

type NewComponentReq struct {
	Data     Data         `json:"data,omitempty"`
	Keyboard Keyboard     `json:"keyboard,omitempty"`
	NextId   *json.Number `json:"next_id,omitempty"`
	Position Point        `json:"position,omitempty"`
}

type Point struct {
	X *json.Number
	Y *json.Number
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
	Text   string       `json:"text,omitempty"`
	NextId *json.Number `json:"next_id,omitempty"`
}
