package userservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"users-service/internal/domain/interfaces/service"
	"users-service/internal/domain/profiles"
	service_error "users-service/internal/service"
	"users-service/pkg/logger/sl"
	umv1 "users-service/proto/gen"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	log         *slog.Logger
	userservice service.IUserService
	umv1.UnimplementedUsersServiceServer
}

func Register(grpc *grpc.Server, userservice service.IUserService, log *slog.Logger) {
	umv1.RegisterUsersServiceServer(grpc, &serverAPI{
		userservice: userservice,
		log:         log,
	})
}

func (s *serverAPI) GetUsers(ctx context.Context, req *emptypb.Empty) (*umv1.GetUsersResponse, error) {
	const op = "grpc.userservice.GetUser"
	log := s.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	users, err := s.userservice.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("users not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "users not found")
		}

		log.Error("Cannot fetch users", sl.Err(err))
		return nil, status.Error(codes.Internal, "cannot fetch users")
	}

	var res_users = make([]*umv1.User, 0, 5)
	for _, user := range users {
		res_users = append(res_users, profiles.UserToProtoUser(user))
	}

	return &umv1.GetUsersResponse{
		Users: res_users,
	}, nil
}

func (s *serverAPI) GetUserById(ctx context.Context, req *umv1.GetUserByIdRequest) (*umv1.GetUserByIdResponse, error) {
	const op = "grpc.userservice.GetUserById"
	log := s.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	req_id := req.GetId()
	if req_id == "" {
		log.Error("id is required", sl.Err(fmt.Errorf("%s: %s", op, "id is required")))
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := uuid.Parse(req_id)
	if err != nil {
		log.Error("wrong id, must be uuid", sl.Err(fmt.Errorf("%s: %s", op, "wrong id, must be uuid")))
		return nil, status.Error(codes.Internal, "wrong id, must be uuid")
	}

	user, err := s.userservice.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("user doesn't exists", sl.Err(err))
			return nil, status.Error(codes.NotFound, "user doesn't exists")
		}

		log.Error("cannot fetch user by id", sl.Err(err))
		return nil, status.Error(codes.Internal, "cannot fetch user by id")
	}

	return &umv1.GetUserByIdResponse{
		User: profiles.UserToProtoUser(user),
	}, nil
}

func (s *serverAPI) InsertUser(ctx context.Context, req *umv1.InsertRequest) (*umv1.InsertResponse, error) {
	const op = "grpc.userservice.insert"
	log := s.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	req_user := req.GetUser()
	if req_user == nil {
		log.Error("user is required", sl.Err(fmt.Errorf("%s: %s", op, "user is required")))
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	user, err := s.userservice.InsertUser(ctx, profiles.ProtoUserToUser(req_user))
	if err != nil {
		if errors.Is(err, service_error.ErrAlreadyExists) {
			log.Warn("user already exists", sl.Err(err))
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		log.Error("cannot insert user", sl.Err(err))
		return nil, status.Error(codes.Internal, "cannot insert user")
	}

	return &umv1.InsertResponse{
		User: profiles.UserToProtoUser(user),
	}, nil
}

func (s *serverAPI) DeleteUser(ctx context.Context, req *umv1.DeleteResuest) (*umv1.DeleteResponse, error) {
	const op = "grpc.userservice.GetUserById"
	log := s.log.With(
		op, "op",
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	req_id := req.GetId()
	if req_id == "" {
		log.Error("id is required", sl.Err(fmt.Errorf("%s: %s", op, "id is required")))
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := uuid.Parse(req_id)
	if err != nil {
		log.Error("wrong id, must be uuid", sl.Err(fmt.Errorf("%s: %s", op, "wrong id, must be uuid")))
		return nil, status.Error(codes.Internal, "wrong id, must be uuid")
	}

	user, err := s.userservice.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("user doesn't exists", sl.Err(err))
			return nil, status.Error(codes.NotFound, "user doesn't exists")
		}

		log.Error("cannot delete user", sl.Err(err))
		return nil, status.Error(codes.Internal, "cannot delete user")
	}

	return &umv1.DeleteResponse{
		User: profiles.UserToProtoUser(user),
	}, nil
}
