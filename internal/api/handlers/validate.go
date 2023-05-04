package handlers

import (
	"errors"
	"regexp"
)

const (
	tokenRegexp = `^\d{9,10}:[\w-]{35}$` //nolint:gosec
)

func validateToken(token string) bool {
	reg := regexp.MustCompile(tokenRegexp)
	return reg.MatchString(token)
}

// Check bot component request struct required fields
func validateAddBotComponent(c *addComponentReq) error {
	if c.Data == nil {
		return errors.New("data.data not found")
	}

	if c.Data.Type == nil {
		return errors.New("data.type not found")
	}

	if c.Data.Content == nil {
		return errors.New("data.content not found")
	}

	// if c.Data.Content.Text == nil {
	// 	return errors.New("data.content.text not found")
	// }

	if c.Commands == nil {
		return errors.New("data.commands not found")
	}

	for _, v := range c.Commands {
		if v.Type == nil {
			return errors.New("data.commands._.type not found")
		}

		if v.Data == nil {
			return errors.New("data.commands._.data not found")
		}
	}

	if c.Position == nil {
		return errors.New("data.Position not found")
	}

	if c.Position.X == nil || c.Position.Y == nil {
		return errors.New("data.Position (x | y param) not found")
	}

	return nil
}

func validateAddCommand(c *addCommandReq) error {
	if c.Type == nil {
		return errors.New("type not found")
	}

	if c.Data == nil {
		return errors.New("data not found")
	}

	if *c.Data == "" {
		return errors.New("data not found")
	}

	return nil
}
