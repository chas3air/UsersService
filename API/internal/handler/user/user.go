package userhandler

import (
	"api/internal/domain/interfaces/service"
	"api/internal/domain/models"
	service_error "api/internal/service"
	"api/pkg/logger/sl"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	log     *slog.Logger
	service service.IUserService
}

func New(log *slog.Logger, service service.IUserService) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}

func (u *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.GetUsersHandler"
	log := u.log.With(
		"op", op,
	)

	users, err := u.service.GetUsers(r.Context())
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("users not found", sl.Err(err))
			WriteUsersToBody(w, http.StatusNotFound, []models.User{})
			return
		}

		log.Error("cannot fetch users", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteUsersToBody(w, http.StatusOK, users)
}

func (u *UserHandler) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.GetUserByIdHandler"
	log := u.log.With(
		"op", op,
	)

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Error("id is required", sl.Err(fmt.Errorf("id is required")))
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	uuid_id, err := uuid.Parse(id)
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	user, err := u.service.GetUserById(r.Context(), uuid_id)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("user not found", sl.Err(err))
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		log.Error("cannot get user by id", sl.Err(err))
		http.Error(w, "cannot get user by id", http.StatusInternalServerError)
		return
	}

	WriteUsersToBody(w, http.StatusOK, user)
}

func (u *UserHandler) InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.InsertUserHandler"
	log := u.log.With(
		"op", op,
	)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error("cannot read and parse request body", sl.Err(err))
		http.Error(w, "cannot read and parse request body", http.StatusBadRequest)
		return
	}

	user, err := u.service.InsertUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, service_error.ErrAlreadyExists) {
			log.Warn("user already exists", sl.Err(err))
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		log.Error("cannot insert user", sl.Err(err))
		http.Error(w, "cannot insert user", http.StatusInternalServerError)
		return
	}

	WriteUsersToBody(w, http.StatusCreated, user)
}

func (u *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.UpdateUserHandler"
	log := u.log.With(
		"op", op,
	)

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Error("id is required", sl.Err(fmt.Errorf("id is required")))
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	uuid_id, err := uuid.Parse(id)
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error("cannot read and parse request body", sl.Err(err))
		http.Error(w, "cannot read and parse request body", http.StatusBadRequest)
		return
	}

	user, err = u.service.UpdateUser(r.Context(), uuid_id, user)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("user not found", sl.Err(err))
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		log.Error("cannot update user", sl.Err(err))
		http.Error(w, "cannot update user", http.StatusInternalServerError)
		return
	}

	WriteUsersToBody(w, http.StatusOK, user)
}

func (u *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.DeleteUserHandler"
	log := u.log.With(
		"op", op,
	)

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Error("id is required", sl.Err(fmt.Errorf("id is required")))
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	uuid_id, err := uuid.Parse(id)
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	user, err := u.service.DeleteUser(r.Context(), uuid_id)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("user not found", sl.Err(err))
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		log.Error("cannot delete user", sl.Err(err))
		http.Error(w, "cannot delete user", http.StatusInternalServerError)
		return
	}

	WriteUsersToBody(w, http.StatusOK, user)
}

func WriteUsersToBody(w http.ResponseWriter, status int, users any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
