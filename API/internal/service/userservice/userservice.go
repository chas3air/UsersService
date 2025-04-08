package userservice

import (
	"api/internal/domain/interfaces/storage"
	"api/internal/domain/models"
	"context"
	"log/slog"

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

// GetUsers implements service.IUserService.
func (u *UserService) GetUsers(context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements service.IUserService.
func (u *UserService) GetUserById(context.Context, uuid.UUID) (models.User, error) {
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
