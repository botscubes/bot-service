package handlers

import (
	"regexp"
)

const (
	tokenRegexp = `^\d{9,10}:[\w-]{35}$` //nolint:gosec
)

func validateToken(token string) bool {
	reg := regexp.MustCompile(tokenRegexp)
	return reg.MatchString(token)
}
