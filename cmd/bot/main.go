package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/botscubes/bot-service/pkg/logger"

	a "github.com/botscubes/bot-service/internal/app"
	"github.com/botscubes/bot-service/internal/config"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Get config: %v\n", err)
		return
	}

	log, err := logger.NewLogger(logger.Config{
		Type: c.LoggerType,
	})
	if err != nil {
		fmt.Printf("Create logger: %v\n", err)
		return
	}

	defer func() {
		if err := log.Sync(); err != nil {
			log.Error(err)
		}
	}()

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := a.CreateApp(log, c)

	go func() {
		<-sigs
		log.Info("Stopping...")

		err = app.Shutdown()
		if err != nil {
			log.Fatalw("Shutdown", "error", err)
		}

		done <- struct{}{}
	}()

	app.Run()

	log.Info("App started")

	<-done

	log.Info("App Done")
}
