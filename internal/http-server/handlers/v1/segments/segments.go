package segments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AlexZahvatkin/segments-users-service/internal/http-server"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/AlexZahvatkin/segments-users-service/internal/use-cases/segments"
	usecases_user_segments "github.com/AlexZahvatkin/segments-users-service/internal/use-cases/user_segments"
	"github.com/go-playground/validator/v10"
)

const (
	maxDescriptionLength = 65536
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=UserAdder
type SegmentAdder interface {
	AddSegment(context.Context, models.AddSegmentParams) (models.Segment, error)
	GetSegmentByName(ctx context.Context, name string) (models.Segment, error)
}

type SegmentDeleter interface {
	DeleteSegment(context.Context, string) error
	GetSegmentByName(ctx context.Context, name string) (models.Segment, error)
}

type SegmentAutoAssigner interface {
	SegmentAdder
	AutoAssigner
}

type AutoAssigner interface {
	GetAllUsersId(ctx context.Context) ([]int64, error)
	AddUserIntoSegment(ctx context.Context, arg models.AddUserIntoSegmentParams) (models.UsersInSegment, error)
}

type responseSegment struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

type responseSegmentAndUsers struct {
	Segment responseSegment `json:"segment"`
	Users   []int64         `json:"added_users_ids"`
}

func AddSegmentHandler(log *slog.Logger, segmentAdder SegmentAutoAssigner) http.HandlerFunc {
	type request struct {
		Name        string  `json:"name" validate:"required,min=4,max=255"`
		Description string  `json:"description"`
		Percent     float64 `json:"percent"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.AddSegmentHandler"

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

		if err := checkDescriptionLength(log, req.Description, w); err != nil {
			return
		}

		req.Name = usecases_segments.FormatSegmnetName(req.Name)

		if _, err := segmentAdder.GetSegmentByName(r.Context(), req.Name); err == nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Segment with such name already exists", log)
			return
		}

		addedSegment, err := segmentAdder.AddSegment(r.Context(), models.AddSegmentParams{
			Name: req.Name,
			Description: sql.NullString{
				String: req.Description,
				Valid:  true,
			},
		})
		if err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not create segment", log)
			return
		}

		respSegm := responseSegment{
			Name:        addedSegment.Name,
			Description: addedSegment.Description.String,
			Created_At:  addedSegment.CreatedAt,
			Updated_At:  addedSegment.UpdatedAt,
		}

		if req.Percent == 0 {
			httpserver.RespondWithJSON(w, http.StatusOK, log, respSegm)
			return
		}

		userIds, err := assignProcentOfUsersToSegment(log, segmentAdder, w, r, req.Percent, req.Name)
		if err != nil {
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, responseSegmentAndUsers{
			Segment: respSegm,
			Users:   userIds,
		})
	}
}

func DeleteSegmentHandler(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.DeleteSegmentHandler"

		handlers.SetLogger(log, r.Context(), op)

		req := r.URL.Query().Get("name")
		if req == "" {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Name is required", log)
			return
		}

		log.Info("Get requested segment name", slog.String("name", req))

		if _, err := segmentDeleter.GetSegmentByName(r.Context(), req); err != nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Invalid segment name", log)
			return
		}

		if err := segmentDeleter.DeleteSegment(r.Context(), req); err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not delete segment:", log)
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, struct{}{})
	}
}

func assignProcentOfUsersToSegment(log *slog.Logger, autoAssigner AutoAssigner, w http.ResponseWriter,
	r *http.Request, percent float64, segmentName string) ([]int64, error) {
	ids, err := autoAssigner.GetAllUsersId(r.Context())
	if err != nil {
		log.Error(err.Error())

		httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not get users", log)
		return nil, err
	}

	pickedIds, err := usecases_user_segments.PickRandomIds(percent, ids)
	if err != nil {
		httpserver.RespondWithError(w, http.StatusBadRequest, err.Error(), log)
		return nil, err
	}

	var res []int64
	for _, id := range pickedIds {
		_, err := autoAssigner.AddUserIntoSegment(r.Context(), models.AddUserIntoSegmentParams{
			UserID:      id,
			SegmentName: segmentName,
		})
		if err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not add user %d", id), log)
			return nil, err
		}
		res = append(res, id)
	}

	return res, nil
}

func checkDescriptionLength(log *slog.Logger, description string, w http.ResponseWriter) error {
	if len(description) > maxDescriptionLength {
		httpserver.RespondWithError(w, http.StatusBadRequest, "Description is too long", log)
		return errors.New("Too long description")
	}
	return nil
}
