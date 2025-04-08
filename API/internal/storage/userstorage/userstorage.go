package userstorage

import (
	"api/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type GRPCUserServer struct {
	log  *slog.Logger
	host string
	port int
}

func New(log *slog.Logger, host string, port int) *GRPCUserServer {
	return &GRPCUserServer{
		log:  log,
		host: host,
		port: port,
	}
}

// GetUsers implements storage.IUserStorage.
func (g *GRPCUserServer) GetUsers(context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements storage.IUserStorage.
func (g *GRPCUserServer) GetUserById(context.Context, uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// InsertUser implements storage.IUserStorage.
func (g *GRPCUserServer) InsertUser(context.Context, models.User) (models.User, error) {
	panic("unimplemented")
}

// DeleteUser implements storage.IUserStorage.
func (g *GRPCUserServer) DeleteUser(context.Context, uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
