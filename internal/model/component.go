package model

import (
	"github.com/botscubes/bot-components/components"
	"github.com/goccy/go-json"

	"github.com/jackc/pgx/v5/pgtype"
)

type ComponentStatus int

var (
	StatusComponentActive ComponentStatus
	StatusComponentDel    ComponentStatus = 1
)

type Component struct {
	components.ComponentData
	Id       int64  `json:"id"`
	Position *Point `json:"position"`
}

type ComponentData struct {
	Type    *string     `json:"type"`
	Content *[]*Content `json:"content"`
}

type Content struct {
	Text *string `json:"text,omitempty"`
}

type Keyboard struct {
	Buttons [][]*int64 `json:"buttons"`
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

// Encode point struct to pgx point type
func (p Point) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: p.X, Y: p.Y},
		Valid: true,
	}, nil
}

// Encode component struct to binary format (for redis)
func (c *Component) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

// Decode component from binary format to struct (fo redis)
func (c *Component) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}

type AddComponentReq struct {
	Data     *ComponentData `json:"data"`
	Commands *CommandsParam `json:"commands"`
	Position *Point         `json:"position"`
}

type SetNextStepComponentReq struct {
	NextStepId *int64 `json:"nextStepId"`
}

type DelSetComponentsReq struct {
	Data *[]int64 `json:"data"`
}

type UpdComponentReq struct {
	Data     *ComponentData `json:"data"`
	Position *Point         `json:"position"`
}
