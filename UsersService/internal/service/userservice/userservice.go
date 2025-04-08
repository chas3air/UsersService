package userservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"users-service/internal/domain/interfaces/storage"
	"users-service/internal/domain/models"
	service_error "users-service/internal/service"
	storage_error "users-service/internal/storage"
	"users-service/pkg/logger/sl"

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
func (u *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "service.GetUsers"
	log := u.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := u.storage.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("users not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("cannot fetch users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// GetUserById implements service.IUserService.
func (u *UserService) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "service.GetUserById"
	log := u.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("user doesn't exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("cannot fetch user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// InsertUser implements service.IUserService.
func (u *UserService) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "service.InsertUser"
	log := u.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.InsertUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage_error.ErrAlreadyExists) {
			log.Warn("user already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, service_error.ErrAlreadyExists)
		}

		log.Error("cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// DeleteUser implements service.IUserService.
func (u *UserService) DeleteUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "service.DeleteUser"
	log := u.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("user not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("cannot delete user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
