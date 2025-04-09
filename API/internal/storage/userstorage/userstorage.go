package userstorage

import (
	"api/internal/domain/models"
	"api/internal/domain/profiles"
	storage_error "api/internal/storage"
	"api/pkg/logger/sl"
	umv1 "api/proto/gen"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
func (g *GRPCUserServer) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.user.GetUsers"
	log := g.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.host, g.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersServiceClient(conn)
	res, err := c.GetUsers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, g.handleError(err, op)
	}

	var res_users = make([]models.User, 0, 5)
	for _, user := range res.GetUsers() {
		res_users = append(res_users, profiles.ProtoUserToUser(user))
	}

	return res_users, nil
}

// GetUserById implements storage.IUserStorage.
func (g *GRPCUserServer) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.user.GetUserById"
	log := g.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.host, g.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersServiceClient(conn)
	res, err := c.GetUserById(ctx, &umv1.GetUserByIdRequest{
		Id: id.String(),
	})
	if err != nil {
		return models.User{}, g.handleError(err, op)
	}

	return profiles.ProtoUserToUser(res.GetUser()), nil
}

// InsertUser implements storage.IUserStorage.
func (g *GRPCUserServer) InsertUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.user.InsertUser"
	log := g.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.host, g.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersServiceClient(conn)
	res, err := c.InsertUser(ctx, &umv1.InsertRequest{
		User: profiles.UserToProtoUser(user),
	})
	if err != nil {
		return models.User{}, g.handleError(err, op)
	}

	return profiles.ProtoUserToUser(res.GetUser()), nil
}

func (g *GRPCUserServer) UpdateUser(ctx context.Context, id uuid.UUID, user models.User) (models.User, error) {
	const op = "storage.user.UpdateUser"
	log := g.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.host, g.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersServiceClient(conn)
	res, err := c.UpdateUser(ctx,
		&umv1.UpdateRequest{
			Id:   id.String(),
			User: profiles.UserToProtoUser(user),
		})
	if err != nil {
		return models.User{}, g.handleError(err, op)
	}

	return profiles.ProtoUserToUser(res.GetUser()), nil
}

// DeleteUser implements storage.IUserStorage.
func (g *GRPCUserServer) DeleteUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.user.DeleteUser"
	log := g.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("request time out", sl.Err(ctx.Err()))
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.host, g.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersServiceClient(conn)
	res, err := c.DeleteUser(ctx, &umv1.DeleteResuest{
		Id: id.String(),
	})
	if err != nil {
		return models.User{}, g.handleError(err, op)
	}

	return profiles.ProtoUserToUser(res.GetUser()), nil
}

func (g *GRPCUserServer) handleError(err error, operation string) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			g.log.Warn("users not found", sl.Err(err))
			return fmt.Errorf("%s: %w", operation, storage_error.ErrNotFound)
		case codes.AlreadyExists:
			g.log.Warn("users already exists", sl.Err(err))
			return fmt.Errorf("%s: %w", operation, storage_error.ErrAlreadyExists)
		default:
			g.log.Error("gRPC error occurred", sl.Err(err))
			return fmt.Errorf("%s: %w", operation, err)
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
