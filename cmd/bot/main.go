package main

import (
	"fmt"

	"github.com/botscubes/bot-service/pkg/logger"

	a "github.com/botscubes/bot-service/internal/app"
	"github.com/botscubes/bot-service/internal/config"
)

func main() {
	var app a.App

	c, err := config.GetConfig()
	if err != nil {
		fmt.Println("Get config: ", err)
		return
	}

	log, err := logger.NewLogger(logger.Config{
		Type: c.LoggerType,
	})
	if err != nil {
		fmt.Println("Create logger: ", err)
		return
	}

	defer func() {
		if err := log.Sync(); err != nil {
			log.Error(err)
		}
	}()

	if err := app.Run(log, c); err != nil {
		log.Error("App run:\n", err)
	}

	log.Info("App Done")
}
