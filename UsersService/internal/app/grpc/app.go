package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"users-service/internal/domain/interfaces/service"
	"users-service/internal/grpc/userservice"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, usersservice service.IUserService, port int) *App {
	gRPCServer := grpc.NewServer()

	userservice.Register(gRPCServer, usersservice, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		"op", op,
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With("op", op).
		Info("stoping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
