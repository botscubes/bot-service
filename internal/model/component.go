package model

import (
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

func (p *Point) ScanPoint(v pgtype.Point) error {
	*p = Point{
		X:     float64(v.P.X),
		Y:     float64(v.P.Y),
		Valid: v.Valid,
	}
	return nil
}

func (p Point) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: float64(p.X), Y: float64(p.Y)},
		Valid: true,
	}, nil
}
