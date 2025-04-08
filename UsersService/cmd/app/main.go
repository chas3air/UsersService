package main

import (
	"users-service/pkg/config"
	"users-service/pkg/logger"
)

func main() {
	config := config.MustLoad()

	log := logger.SetupLogger(config.Env)

	log.Info("Ready to change")
}
