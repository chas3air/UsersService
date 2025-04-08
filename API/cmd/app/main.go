package app

import (
	"api/internal/app"
	"api/pkg/config"
	"api/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := config.MustLoad()

	log := logger.SetupLogger(config.Env)

	appplication := app.New(log, config)

	go func() {
		appplication.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("application is stopped")
}
