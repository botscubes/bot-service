package handlers

// TODO: create pkg

import (
	"github.com/botscubes/bot-service/internal/config"
)

func CheckIsMain(id int64) bool {
	return id == config.MainComponentId
}
