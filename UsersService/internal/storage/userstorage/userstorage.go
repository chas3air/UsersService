package userstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"users-service/internal/domain/models"
	storage_error "users-service/internal/storage"
	"users-service/pkg/logger/sl"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type PsqlStorage struct {
	log *slog.Logger
	DB  *sql.DB
}

const UsersTableName = "Users"

func New(log *slog.Logger, connStr string) *PsqlStorage {
	const op = "psql.New"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.With(slog.String("op", op)).Error("Error connecting to DB", sl.Err(err))
		panic(err)
	}

	wd, _ := os.Getwd()
	migrationPath := filepath.Join(wd, "app", "migrations")
	if err := applyMigrations(db, migrationPath); err != nil {
		panic(err)
	}

	return &PsqlStorage{
		log: log,
		DB:  db,
	}
}

func applyMigrations(db *sql.DB, migrationsPath string) error {
	return goose.Up(db, migrationsPath)
}

// GetUsers implements storage.IUserStorage.
func (p *PsqlStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "service.GetUsers"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	rows, err := p.DB.QueryContext(ctx, `
		SELECT * FROM `+UsersTableName+`;
	`)
	if err != nil {
		log.Error("error scanning rows", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users = make([]models.User, 0, 5)
	var buff_user models.User
	for rows.Next() {
		err := rows.Scan(&buff_user.Id, &buff_user.Login, &buff_user.Password)
		if err != nil {
			log.Warn("cannot scan row", sl.Err(err))
			continue
		}

		users = append(users, buff_user)
	}

	if len(users) == 0 {
		log.Warn("no users found")
		return nil, storage_error.ErrNotFound
	}

	return users, nil
}

// GetUserById implements storage.IUserStorage.
func (p *PsqlStorage) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "service.GetUsers"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	err := p.DB.QueryRowContext(ctx, `
		SELECT * FROM `+UsersTableName+`
		WHERE id=$1
	`, id).Scan(&user.Id, &user.Login, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user not found", sl.Err(err))
			return models.User{}, storage_error.ErrNotFound
		}
		log.Error("cannot scan user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// InsertUser implements storage.IUserStorage.
func (p *PsqlStorage) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "service.GetUsers"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	_, err := p.DB.ExecContext(ctx, `
		INSERT INTO `+UsersTableName+`(id, login, password)
		VALUES($1, $2, $3);
	`, user.Id, user.Login, user.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Warn("user already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: user already exists", op)
		}

		log.Error("cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// DeleteUser implements storage.IUserStorage.
func (p *PsqlStorage) DeleteUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "service.GetUsers"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := p.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("user doesn't exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("cannot delete user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.ExecContext(ctx, `
		DELETE FROM `+UsersTableName+`
		WHERE id=$1;
	`, id)
	if err != nil {
		log.Error("cannot delete user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
