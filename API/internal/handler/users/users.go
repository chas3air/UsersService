package usershandler

import (
	"api/internal/domain/interfaces/service"
	"log/slog"
)

type UserHandler struct {
	log     *slog.Logger
	service service.IUserService
}

func New(log *slog.Logger, service service.IUserService) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}
