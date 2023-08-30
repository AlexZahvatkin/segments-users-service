package users_in_segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"errors"

	httpserver "github.com/AlexZahvatkin/segments-users-service/internal/http-server"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
)

const (
	timeFormat = "2006-01-02T15:04:05Z07:00"	//RFC3339 
)

type SegmentsAssigner interface {
	AddUserIntoSegment(context.Context, models.AddUserIntoSegmentParams) (models.UsersInSegment, error)
	RemoveUserFromSegment(ctx context.Context, arg models.RemoveUserFromSegmentParams) error
	SegmentGetter
	UserGetter
}

type SegmentsForUserGetter interface {
	GetSegmentsByUserId(ctx context.Context, userID int64) ([]string, error)
	UserGetter
}

type SegmentHistoryGetter interface {
	GetSegmentsHistoryByUserId(ctx context.Context, arg models.GetSegmentsHistoryByUserIdParams) ([]models.UsersInSegmentsHistory, error)
	UserGetter
}

type SegmentsAssignerWithTTL interface {
	AddUserIntoSegmentWithTTLInHours(context.Context, models.AddUserIntoSegmentWithTTLInHoursParams) (models.UsersInSegment, error)
	GetSegmentByName(ctx context.Context, name string) (models.Segment, error)
	UserGetter
}

type UserGetter interface {
	GetUserById(context.Context, int64) (models.User, error)
}

type SegmentGetter interface {
	GetSegmentByName(ctx context.Context, name string) (models.Segment, error)
}

type UsersInSegmentsResponse struct {
	UserId      int64     `json:"user_id"`
	SegmentName string    `json:"segment_name"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
	Expire_at   time.Time `json:"expire_at,omitempty"`
}

type UsersInSegmentsHistoryResponse struct {
	UserId      int64     `json:"user_id"`
	SegmentName string    `json:"segment_name"`
	ActionType  string    `json:"action_type"`
	ActionDate  time.Time `json:"action_date"`
}

func SegmentsAssignHandler(log *slog.Logger, assigner SegmentsAssigner) http.HandlerFunc {
	type request struct {
		SegmentsToDeleteNames []string `json:"to_delete"`
		SegmentsToAddNames    []string `json:"to_add"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.SegmentsAssignHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		decoder := json.NewDecoder(r.Body)
		req := request{}
		err := decoder.Decode(&req)
		if err != nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err), log)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err!= nil {
			httpserver.RespondWithValidateError(w, log, err)
			return
		}

		userId, err := getUserIdFromParams(w, r, log)
		if err!= nil {
            return
        }

		if !checkIfUserExists(assigner, log, userId, w, r) {
			return
		}

		for _, segmentName := range req.SegmentsToAddNames {
			if !checkIfSegmentExists(assigner, log, segmentName, w, r) {
				return
			}
		}

		for _, segmentName := range req.SegmentsToDeleteNames {
			if !checkIfSegmentExists(assigner, log, segmentName, w, r) {
				return
			}
		}

		for _, segmentName := range req.SegmentsToDeleteNames {
			if err := assigner.RemoveUserFromSegment(r.Context(),
				models.RemoveUserFromSegmentParams{UserID: userId, SegmentName: segmentName}); err != nil {
				log.Error(err.Error())

				httpserver.RespondWithError(w, http.StatusInternalServerError,
					fmt.Sprintf("Failed to delete segment %s for user %d", segmentName, userId), log)
				return
			}
		}

		var result []models.UsersInSegment

		for _, segmentName := range req.SegmentsToAddNames {
			res, err := assigner.AddUserIntoSegment(r.Context(), models.AddUserIntoSegmentParams{UserID: userId, SegmentName: segmentName})
			if err != nil {
				log.Error(err.Error())

				httpserver.RespondWithError(w, http.StatusInternalServerError,
					fmt.Sprintf("Failed to add segment %s for user %d", segmentName, userId), log)
				return
			}
			result = append(result, res)
		}

		var resp []UsersInSegmentsResponse
		for _, r := range result {
			resp = append(resp, transformToUsersInSegmentsResponse(r))
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, resp)
	}
}

func SegmentsAssignWithTTLInHoursHandler(log *slog.Logger, assigner SegmentsAssignerWithTTL) http.HandlerFunc {
	type request struct {
		SegmentName string `json:"segment_name" validate:"required"`
		TTL         int32  `json:"ttl" validate:"required,gt=0"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.SegmentsAssignHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		decoder := json.NewDecoder(r.Body)
		req := request{}
		err := decoder.Decode(&req)
		if err != nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err), log)
			return
		}

		userId, err := getUserIdFromParams(w, r, log)
		if err!= nil {
            return
        }

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err!= nil {
			httpserver.RespondWithValidateError(w, log, err)
			return
		}

		if !checkIfUserExists(assigner, log, userId, w, r) {
			return
		}

		if !checkIfSegmentExists(assigner, log, req.SegmentName, w, r) {
			return
		}

		res, err := assigner.AddUserIntoSegmentWithTTLInHours(r.Context(), models.AddUserIntoSegmentWithTTLInHoursParams{
			UserID: userId,
			SegmentName: req.SegmentName,
            NumberOfHours: req.TTL,
		})
		if err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to add segment %s for user %d", req.SegmentName, userId), log)
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, transformToUsersInSegmentsResponse(res))
	}
}

func GetSegmentsForUserHandler(log *slog.Logger, getter SegmentsForUserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.GetSegmentsForUserHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, err := getUserIdFromParams(w, r, log)
		if err!= nil {
            return
        }

		if !checkIfUserExists(getter, log, userId, w, r) {
			return
		}

		res, err := getter.GetSegmentsByUserId(r.Context(), userId)
		if err != nil {
			if err == sql.ErrNoRows {
				httpserver.RespondWithJSON(w, http.StatusNoContent, log, struct{}{})
				return
			}
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get segments for user %d", userId), log)
			return
		}

		if len(res) == 0 {
			httpserver.RespondWithJSON(w, http.StatusNoContent, log, struct{}{})
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, res)
	}
}

func GetSegmentsHistoryByUser(log *slog.Logger, getter SegmentHistoryGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.GetSegmentsHistoryByUser"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, err := getUserIdFromParams(w, r, log)
		if err!= nil {
            return
        }

		from, err := getTimeFromParams(w, r, log, "from")
		if err!= nil {
			return
		}

		to, err := getTimeFromParams(w, r, log, "to")
        if err!= nil {
            return
        }

		if !checkIfUserExists(getter, log, userId, w, r) {
			return
		}

		res, err := getter.GetSegmentsHistoryByUserId(r.Context(), models.GetSegmentsHistoryByUserIdParams{
			UserID:   userId,
			FromDate: from,
			ToDate:   to,
		})

		if err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get segments for user %d", userId), log)
			return
		}

		var resp [][]string
		for _, item := range res {
			resp = append(resp, transformUsersInSegmentsHistoryToString(item))
		}

		httpserver.RespondWithCSV(w, http.StatusOK, log, resp)
	}
}

func transformToUsersInSegmentsResponse(userInSegment models.UsersInSegment) UsersInSegmentsResponse {
	return UsersInSegmentsResponse{
		UserId:      userInSegment.UserID,
		SegmentName: userInSegment.SegmentName,
		Created_At:  userInSegment.CreatedAt,
		Updated_At:  userInSegment.UpdatedAt,
		Expire_at:   userInSegment.ExpireAt.Time,
	}
}

func transformUsersInSegmentsHistoryToString(userInSegment models.UsersInSegmentsHistory) []string {
	var res []string
	res = append(res, strconv.FormatInt(userInSegment.UserID, 10))
	res = append(res, userInSegment.SegmentName)
	res = append(res, userInSegment.ActionType)
	res = append(res, userInSegment.ActionDate.Format("2006-01-02 15:04:05"))
	return res
}

func checkIfUserExists(getter UserGetter, log *slog.Logger, userId int64, w http.ResponseWriter, r *http.Request) bool {
	if _, err := getter.GetUserById(r.Context(), userId); err != nil {
		if err == sql.ErrNoRows {
			httpserver.RespondWithError(w, http.StatusBadRequest, "User does not exist", log)
			return false
		}
		httpserver.RespondWithError(w, http.StatusInternalServerError, "Failed to get user", log)
		return false
	}
	return true
}

func checkIfSegmentExists(getter SegmentGetter, log *slog.Logger, segment_name string, w http.ResponseWriter, r *http.Request) bool {
	_, err := getter.GetSegmentByName(r.Context(), segment_name)
	if err != nil {
		if err == sql.ErrNoRows {
			httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Segment %s does not exist", segment_name), log)
			return false
		}
		httpserver.RespondWithError(w, http.StatusInternalServerError, "Failed to get segment", log)
		return false
	}
	return true
}

func getUserIdFromParams(w http.ResponseWriter, r *http.Request, log *slog.Logger) (int64, error) {
	s := chi.URLParam(r, "userId")
	log.Debug("Param is " + s)
	if s == "" {
		httpserver.RespondWithError(w, http.StatusBadRequest, "You must provide userId", log)
		return -1, errors.New("No parameter provided")
	}
	userId, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("User id must be a number: %v", err), log)
		return -1, err
	}
	return userId, nil
}

func getTimeFromParams(w http.ResponseWriter, r *http.Request, log *slog.Logger, paramName string) (time.Time, error) {
	s := r.URL.Query().Get(paramName)
	if s == "" {
		httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("You must provide %s", paramName), log)
		return time.Time{}, errors.New("No parameter provided")
	}
	t, err := time.Parse(timeFormat, s)
	if err != nil {		
		httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Wrong datetime format: %v", err), log)
		return time.Time{}, err
	}
	return t, nil
}