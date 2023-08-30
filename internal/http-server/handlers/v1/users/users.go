package users

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexZahvatkin/segments-users-service/internal/http-server"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=UserAdder
type UserAdder interface {
	AddUser(context.Context, string) (models.User, error)
}

// @Summary Add new user
// @Description Creates new user with a given name
// @Tags Users 
// @Accept  json
// @Produce  json
// @ID create-user
// @Param name body string true "User name"
// @Success 201 {object} models.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /v1/users [post]
func AddUserHandler(log *slog.Logger, userAdder UserAdder) http.HandlerFunc {
	type request struct {
		Name string `json:"name" validate:"required,min=4,max=255"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.AddUserHandler"

		handlers.SetLogger(log, r.Context(), op)

		req, err := httpserver.DecodeRequsetBody(w, r, request {}, log)
		if err != nil {
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

		httpserver.RespondWithJSON(w, http.StatusCreated, log, user)
	}
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=UserDeleter
type UserDeleter interface {
	DeleteUser(context.Context, int64) error
	GetUserById(context.Context, int64) (models.User, error)
}

// @Summary Delete user
// @Description Deletes user with a given id
// @Tags Users 
// @Accept  json
// @Produce  json
// @ID delete-user
// @Param id path int true "User ID"
// @Success 200
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /v1/users [delete]
func DeleteUserHandler(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.DeleteUserHandler"

		handlers.SetLogger(log, r.Context(), op)

		userId, err := httpserver.GetUserIdFromParams(w, r, log)
		if err!= nil {
            return
        }

		if _, err := userDeleter.GetUserById(r.Context(), userId); err != nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Invalid user id", log)
			return
		}
		
		if err = userDeleter.DeleteUser(r.Context(), userId); err!= nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not delete user:", log)
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, struct{}{})
	}
}