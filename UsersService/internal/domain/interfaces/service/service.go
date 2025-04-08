package service

import (
	"context"
	"users-service/internal/domain/models"

	"github.com/google/uuid"
)

type IUserService interface {
	GetUsers(context.Context) ([]models.User, error)
	GetUserById(context.Context, uuid.UUID) (models.User, error)
	InsertUser(context.Context, models.User) (models.User, error)
	DeleteUser(context.Context, uuid.UUID) (models.User, error)
}
