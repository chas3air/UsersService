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
	user_storage := userstorage.New(a.log, a.config.ServerHost, a.config.ServerPort)
	user_service := userservice.New(a.log, user_storage)
	user_handler := userhandler.New(a.log, user_service)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/api/v1/users", user_handler.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", user_handler.GetUserByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", user_handler.InsertUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{id}", user_handler.UpdateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/users/{id}", user_handler.DeleteUserHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.config.Api.Port), r); err != nil {
		panic(err)
	}
}
