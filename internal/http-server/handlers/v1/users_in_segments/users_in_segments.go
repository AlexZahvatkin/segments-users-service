package users_in_segments

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	httpserver "github.com/AlexZahvatkin/segments-users-service/internal/http-server"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/go-playground/validator/v10"
)

const (
	timeFormat = "2006-01-02T15:04:05Z07:00" //RFC3339
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

// @Summary Assigns segments to a user.
// @Description Adds and deletes segments provided by a request for user with provied id.
// @Tags Useres in segments
// @Accept  json
// @Produce  json
// @ID segments-assign
// @Param userId path int true "User id"
// @Param segments body models.SegmentAssignRequest true "Segments to delete and add for user"
// @Success 200 {object} UsersInSegmentsResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /v1/segments/assign/{userId} [post]
func SegmentsAssignHandler(log *slog.Logger, assigner SegmentsAssigner) http.HandlerFunc {
	type request struct {
		SegmentsToDeleteNames []string `json:"to_delete"`
		SegmentsToAddNames    []string `json:"to_add"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.SegmentsAssignHandler"

		handlers.SetLogger(log, r.Context(), op)

		req, err := httpserver.DecodeRequsetBody(w, r, request{}, log)
		if err != nil {
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			httpserver.RespondWithValidateError(w, log, err)
			return
		}

		userId, err := httpserver.GetUserIdFromParams(w, r, log)
		if err != nil {
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

// @Summary Assigns segments to a user with ttl.
// @Description Adds a provided segment to a provided user with TTL in hours.
// @Tags Useres in segments
// @Accept  json
// @Produce  json
// @ID segments-assign-with-ttl
// @Param userId path int true "User id"
// @Param segments body models.SegmentAssignWithTTLRequest true "Segment to assign and TTL in hours"
// @Success 200 {object} UsersInSegmentsResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /v1/segments/ttl/{userId} [post]
func SegmentsAssignWithTTLInHoursHandler(log *slog.Logger, assigner SegmentsAssignerWithTTL) http.HandlerFunc {
	type request struct {
		SegmentName string `json:"segment_name" validate:"required"`
		TTL         int32  `json:"ttl" validate:"required,gt=0"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.SegmentsAssignHandler"

		handlers.SetLogger(log, r.Context(), op)

		req, err := httpserver.DecodeRequsetBody(w, r, request{}, log)
		if err != nil {
			return
		}

		userId, err := httpserver.GetUserIdFromParams(w, r, log)
		if err != nil {
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
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
			UserID:        userId,
			SegmentName:   req.SegmentName,
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

// @Summary Segments for user
// @Description Returns a list of segments that are active for a provided user.
// @Tags Useres in segments
// @Accept  json
// @Produce  json
// @ID get-segments-for-user
// @Param userId path int true "User id"
// @Success 200 {object} []string
// @Success 204
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /v1/segments/{userId} [get]
func GetSegmentsForUserHandler(log *slog.Logger, getter SegmentsForUserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.GetSegmentsForUserHandler"

		handlers.SetLogger(log, r.Context(), op)

		userId, err := httpserver.GetUserIdFromParams(w, r, log)
		if err != nil {
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

// @Summary Segments history for user
// @Description Returns a history of added and deleted segments for a provided user in a given period.
// @Tags Useres in segments
// @Accept  json
// @Produce  json
// @ID get-segments-for-user-history
// @Param userId path int true "User id"
// @Param from path string true "From datetime"
// @Param to path string true "To datetime"
// @Success 200 {object} []string
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /v1/segments/history/{userId} [get]
func GetSegmentsHistoryByUser(log *slog.Logger, getter SegmentHistoryGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.GetSegmentsHistoryByUser"

		handlers.SetLogger(log, r.Context(), op)

		userId, err := httpserver.GetUserIdFromParams(w, r, log)
		if err != nil {
			return
		}

		from, err := httpserver.GetTimeFromParams(w, r, log, "from", timeFormat)
		if err != nil {
			return
		}

		to, err := httpserver.GetTimeFromParams(w, r, log, "to", timeFormat)
		if err != nil {
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
