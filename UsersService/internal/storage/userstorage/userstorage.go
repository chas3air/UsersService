package userstorage

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"users-service/internal/domain/models"
	"users-service/pkg/logger/sl"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
func (p *PsqlStorage) GetUsers(context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements storage.IUserStorage.
func (p *PsqlStorage) GetUserById(context.Context, uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// InsertUser implements storage.IUserStorage.
func (p *PsqlStorage) InsertUser(context.Context, models.User) (models.User, error) {
	panic("unimplemented")
}

// DeleteUser implements storage.IUserStorage.
func (p *PsqlStorage) DeleteUser(context.Context, uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
