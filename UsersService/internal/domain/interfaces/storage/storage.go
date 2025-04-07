package storage

import (
	"context"
	"users-service/internal/domain/models"

	"github.com/google/uuid"
)

type IUserStorage interface {
	GetUsers(context.Context) ([]models.User, error)
	GetUserById(context.Context, string) (models.User, error)
	InsertUser(context.Context, models.User) (models.User, error)
	DeleteUser(context.Context, uuid.UUID) (models.User, error)
}
