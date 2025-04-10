package app

import (
	"log/slog"
	grpcapp "users-service/internal/app/grpc"
	"users-service/internal/domain/interfaces/storage"
	"users-service/internal/service/userservice"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, storage storage.IUserStorage, port int) *App {
	userService := userservice.New(log, storage)

	grpcApp := grpcapp.New(log, userService, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
