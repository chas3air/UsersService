package main

import (
	"os"
	"os/signal"
	"syscall"
	"users-service/internal/app"
	"users-service/pkg/config"
	"users-service/pkg/logger"
)

func main() {
	config := config.MustLoad()

	log := logger.SetupLogger(config.Env)

	log.Info("Ready to change")

	application := app.New(log, nil, config.Grpc.Port)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("Stoping application")
	application.GRPCServer.Stop()

	log.Info("Stoping db")
	// storage.Close()

	log.Info("application is stopped")
}
