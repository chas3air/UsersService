package storage

import (
	"api/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type IUserStorage interface {
	GetUsers(context.Context) ([]models.User, error)
	GetUserById(context.Context, uuid.UUID) (models.User, error)
	InsertUser(context.Context, models.User) (models.User, error)
	UpdateUser(context.Context, uuid.UUID, models.User) (models.User, error)
	DeleteUser(context.Context, uuid.UUID) (models.User, error)
}
