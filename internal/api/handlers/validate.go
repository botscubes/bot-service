package handlers

import (
	"errors"
)

// Check component required fields
func validateComponent(c *addComponentReq) error {
	if c.Data == nil {
		return errors.New("data.data not found")
	} else {
		if c.Data.Type == nil {
			return errors.New("data.type not found")
		}

		if c.Data.Content == nil {
			return errors.New("data.content not found")
		} else {
			if c.Data.Content.Text == nil {
				return errors.New("data.content.text not found")
			}
		}
	}

	if c.Commands == nil {
		return errors.New("data.commands not found")
	} else {
		for _, v := range c.Commands {
			if v.Type == nil {
				return errors.New("data.commands._.type not found")
			}

			if v.Data == nil {
				return errors.New("data.commands._.data not found")
			}
		}
	}

	if c.Position == nil {
		return errors.New("data.Position not found")
	} else {
		if c.Position.X == nil || c.Position.Y == nil {
			return errors.New("data.Position (x | y param) not found")
		}
	}

	return nil
}
