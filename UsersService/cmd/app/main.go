package main

import (
	"users-service/pkg/config"
	"users-service/pkg/logger"
)

func main() {
	config := config.MustLoad()

	logger := logger.SetupLogger(config.Env)

	
}
