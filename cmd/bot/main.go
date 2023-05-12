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

	defer app.Log.Sync()

	if err := app.Run(log); err != nil {
		log.Fatal("App run:\n", err)
	}

	log.Info("App Done")
}
