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
	storageerror "users-service/internal/storage"
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

func (p *PsqlStorage) Close() {
	p.DB.Close()
}

// GetUsers implements storage.IUserStorage.
func (p *PsqlStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.user.GetUsers"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	rows, err := p.DB.QueryContext(ctx, `
		SELECT * FROM `+UsersTableName+`;
	`)
	if err != nil {
		log.Error("error scanning rows", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, storageerror.ErrNotFound)
	}
	defer rows.Close()

	var users = make([]models.User, 0, 5)
	var buffUser models.User
	for rows.Next() {
		err := rows.Scan(&buffUser.Id, &buffUser.Login, &buffUser.Password)
		if err != nil {
			log.Warn("cannot scan row", sl.Err(err))
			continue
		}

		users = append(users, buffUser)
	}

	return users, nil
}

// GetUserById implements storage.IUserStorage.
func (p *PsqlStorage) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.user.GetUserById"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
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
			return models.User{}, storageerror.ErrNotFound
		}

		log.Error("cannot scan user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// InsertUser implements storage.IUserStorage.
func (p *PsqlStorage) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.user.InsertUser"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
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
			return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrAlreadyExists)
		}

		log.Error("cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (p *PsqlStorage) UpdateUser(ctx context.Context, id uuid.UUID, user models.User) (models.User, error) {
	const op = "storage.user.UpdateUser"
	log := p.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	result, err := p.DB.ExecContext(ctx, `
		UPDATE `+UsersTableName+`
		SET login=$1, password=$2
		WHERE id=$3
	`, user.Login, user.Password, id)
	if err != nil {
		log.Error("fialed to update user", sl.Err(err))
		return user, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Error get rows affected", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		log.Error("Zero rows affected")
		return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrNotFound)
	}

	return user, nil
}

// DeleteUser implements storage.IUserStorage.
func (p *PsqlStorage) DeleteUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.user.DeleteUser"
	log := p.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := p.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, storageerror.ErrNotFound) {
			log.Warn("user doesn't exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrNotFound)
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
