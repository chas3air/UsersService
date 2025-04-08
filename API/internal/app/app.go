package app

import (
	usershandler "api/internal/handler/users"
	"api/internal/service/userservice"
	"api/internal/storage/userstorage"
	"api/pkg/config"
	"log/slog"
)

type App struct {
	log    *slog.Logger
	config *config.Config
}

func New(log *slog.Logger, config *config.Config) *App {
	return &App{}
}

func (a *App) Run() {
	user_storage := userstorage.New(a.log, a.config.ServerHost, a.config.ServerPort)
	user_service := userservice.New(a.log, user_storage)
	user_handler := usershandler.New(a.log, user_service)

	
}
