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
	userservice := userservice.New(log, storage)

	grpcapp := grpcapp.New(log, userservice, port)

	return &App{
		GRPCServer: grpcapp,
	}
}
