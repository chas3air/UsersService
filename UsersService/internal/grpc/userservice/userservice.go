package userservice

import (
	"context"
	"log/slog"
	"users-service/internal/domain/interfaces/service"
	umv1 "users-service/proto/gen"

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

	return nil, nil
}

func (s *serverAPI) GetUserById(context.Context, *umv1.GetUserByIdRequest) (*umv1.GetUserByIdResponse, error) {
	const op = "grpc.userservice.GetUserById"
	log := s.log.With(
		op, "op",
	)

	return nil, nil
}

func (s *serverAPI) InsertUser(context.Context, *umv1.InsertRequest) (*umv1.InsertResponse, error) {
	const op = "grpc.userservice.insert"
	log := s.log.With(
		op, "op",
	)

	return nil, nil
}

func (s *serverAPI) DeleteUser(context.Context, *umv1.DeleteResuest) (*umv1.DeleteResponse, error) {
	const op = "grpc.userservice.GetUserById"
	log := s.log.With(
		op, "op",
	)

	return nil, nil
}
