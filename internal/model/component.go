package model

import (
	"github.com/goccy/go-json"

	"github.com/jackc/pgx/v5/pgtype"
)

type Component struct {
	Id         int64     `json:"id"`
	Data       *Data     `json:"data"`
	Keyboard   *Keyboard `json:"keyboard"`
	Commands   *Commands `json:"commands"`
	NextStepId *int64    `json:"next_step_id"`
	IsMain     bool      `json:"is_main"`
	Position   *Point    `json:"position"`
	Status     int       `json:"-"`
}

type Commands []*Command

type Data struct {
	Type    *string     `json:"type"`
	Content *[]*Content `json:"content"`
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
	ComponentId *int64  `json:"component_id"`
	NextStepId  *int64  `json:"next_step_id"`
	Status      int     `json:"status"`
}

type Point struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Valid bool    `json:"-"`
}

// Decode pgx point type to point struct
func (p *Point) ScanPoint(v pgtype.Point) error {
	*p = Point{
		X:     v.P.X,
		Y:     v.P.Y,
		Valid: v.Valid,
	}
	return nil
}

// Encode point strcut to pgx point type
func (p Point) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: p.X, Y: p.Y},
		Valid: true,
	}, nil
}

// Mb rewrite to https://github.com/nitishm/go-rejson

// Encode component struct to binary format (for redis)
func (c *Component) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

// Decode component from binary format to struct (fo redis)
func (c *Component) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}
