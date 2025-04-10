package app

import (
	userhandler "api/internal/handler/user"
	"api/internal/service/userservice"
	"api/internal/storage/userstorage"
	"api/pkg/config"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	log    *slog.Logger
	config *config.Config
}

func New(log *slog.Logger, config *config.Config) *App {
	return &App{
		log:    log,
		config: config,
	}
}

func (a *App) Run() {
	userStorage := userstorage.New(a.log, a.config.ServerHost, a.config.ServerPort)
	userService := userservice.New(a.log, userStorage)
	userHandler := userhandler.New(a.log, userService)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/api/v1/users", userHandler.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", userHandler.GetUserByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", userHandler.InsertUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{id}", userHandler.UpdateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/users/{id}", userHandler.DeleteUserHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.config.Api.Port), r); err != nil {
		panic(err)
	}
}
