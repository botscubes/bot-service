package main

import (
	"fmt"

	"github.com/botscubes/bot-service/pkg/logger"

	a "github.com/botscubes/bot-service/internal/app"
)

func main() {
	var app a.App

	log, err := logger.NewLogger()
	if err != nil {
		_ = fmt.Errorf("new logger: %w", err)
	}

	defer func() {
		if err := app.Log.Sync(); err != nil {
			log.Error(err)
		}
	}()

	if err := app.Run(log); err != nil {
		log.Error("App run:\n", err)
	}

	log.Info("App Done")
}
