package main

import (
	"api/internal/app"
	"api/pkg/config"
	"api/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := config.MustLoad()

	log := logger.SetupLogger(config.Env)

	log.Info("starting application", slog.Any("config", config))

	appplication := app.New(log, config)

	go func() {
		appplication.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("application is stopped")
}
