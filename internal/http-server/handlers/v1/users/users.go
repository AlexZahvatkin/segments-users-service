package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/AlexZahvatkin/segments-users-service/internal/http-server"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=UserAdder
type UserAdder interface {
	AddUser(context.Context, string) (models.User, error)
}

func AddUserHandler(log *slog.Logger, userAdder UserAdder) http.HandlerFunc {
	type request struct {
		Name string `json:"name" validate:"required,min=4,max=255"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.AddUserHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		decoder := json.NewDecoder(r.Body) 
		req := request {}
		err := decoder.Decode(&req)
		if err!= nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err), log)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err!= nil {
			httpserver.RespondWithValidateError(w, log, err)
			return
		}

		user, err := userAdder.AddUser(r.Context(), req.Name)
		if err != nil {
			log.Error(err.Error())
		
			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not create user:", log)
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, user)
	}
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=UserDeleter
type UserDeleter interface {
	DeleteUser(context.Context, int64) error
	GetUserById(context.Context, int64) (models.User, error)
}

func DeleteUserHandler(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	type request struct {
		ID int64 `json:"id" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.DeleteUserHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		decoder := json.NewDecoder(r.Body) 
		req := request {}
		err := decoder.Decode(&req)
		if err!= nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err), log)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err!= nil {
			httpserver.RespondWithValidateError(w, log, err)
			return
		}

		if _, err := userDeleter.GetUserById(r.Context(), req.ID); err != nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Invalid user id", log)
			return
		}
		
		if err = userDeleter.DeleteUser(r.Context(), req.ID); err!= nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not delete user:", log)
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, struct{}{})
	}
}