package main

import (
	"github.com/botscubes/bot-service/pkg/log"

	a "github.com/botscubes/bot-service/internal/app"
)

func main() {
	var app a.App
	if err := app.Run(); err != nil {
		log.Fatal("App run:\n", err)
	}

	log.Info("App Done")
}
