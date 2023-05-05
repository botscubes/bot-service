package components

import (
	"errors"
)

var (
	mainComponentId int64 = 1
	maxPositionX    int64 = 10000
	maxPositionY    int64 = 10000
)

func CheckIsMain(id int64) bool {
	return id == mainComponentId
}

func (ct *Components) ValidateData(d *Data) error {
	if d == nil {
		return errors.New("data not found")
	}

	if d.Type == nil {
		return errors.New("data.type not found")
	}

	if d.Content == nil {
		return errors.New("data.content not found")
	}

	switch *d.Type {
	case "start":
		return vStart(d.Content)
	case "text":
		return vText(d.Content)
	default:
		return errors.New("unknown data.type")
	}
}

func (ct *Components) ValidateCommands(c *[]*Command) error {
	if c == nil {
		return errors.New("commands not found")
	}

	for _, v := range *c {
		if err := ct.ValidateCommand(v.Type, v.Data); err != nil {
			return err
		}
	}

	return nil
}

func (ct *Components) ValidateCommand(t, d *string) error {
	if t == nil {
		return errors.New("type not found")
	}

	if d == nil {
		return errors.New("data not found")
	}

	if *d == "" {
		return errors.New("data not found")
	}

	return nil
}

func ValidatePosition(p *Point) error {
	if p == nil {
		return errors.New("Position not found")
	}

	if p.X == nil || p.Y == nil {
		return errors.New("Position: x or y param not found")
	}

	if int64(*p.X) < 0 || int64(*p.X) > maxPositionX {
		return errors.New("Position: x param has an incorrect value")
	}

	if int64(*p.Y) < 0 || int64(*p.Y) > maxPositionY {
		return errors.New("Position: x param has an incorrect value")
	}

	return nil
}

func vStart(c *[]*Content) error {
	return errors.New("its start component")
}

func vText(c *[]*Content) error {
	if len(*c) != 1 {
		return errors.New("invalid content len")
	}

	return nil
}
