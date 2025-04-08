package userservice

import (
	"context"
	"log/slog"
	"users-service/internal/domain/interfaces/storage"
	"users-service/internal/domain/models"

	"github.com/google/uuid"
)

type UserService struct {
	log     *slog.Logger
	storage storage.IUserStorage
}

func New(log *slog.Logger, storage storage.IUserStorage) *UserService {
	return &UserService{
		log:     log,
		storage: storage,
	}
}

// GetUserById implements service.IUserService.
func (u *UserService) GetUserById(context.Context, string) (models.User, error) {
	panic("unimplemented")
}

// GetUsers implements service.IUserService.
func (u *UserService) GetUsers(context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// InsertUser implements service.IUserService.
func (u *UserService) InsertUser(context.Context, models.User) (models.User, error) {
	panic("unimplemented")
}

// DeleteUser implements service.IUserService.
func (u *UserService) DeleteUser(context.Context, uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
