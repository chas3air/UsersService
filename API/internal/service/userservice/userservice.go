package userservice

import (
	"api/internal/domain/interfaces/storage"
	"api/internal/domain/models"
	service_error "api/internal/service"
	storage_error "api/internal/storage"
	"api/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
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
func (u *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "service.user.GetUsers"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := u.storage.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("users not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("caannot fetch users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
	}

	return users, nil
}

// GetUserById implements service.IUserService.
func (u *UserService) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "service.user.GetUserById"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("user doesn't exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("cannot get user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// InsertUser implements service.IUserService.
func (u *UserService) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "service.user.InsertUser"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
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
	const op = "service.user.DeleteUser"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
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
