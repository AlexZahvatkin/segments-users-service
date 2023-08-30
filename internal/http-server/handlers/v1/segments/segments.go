package segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AlexZahvatkin/segments-users-service/internal/http-server"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/AlexZahvatkin/segments-users-service/internal/use-cases/segments"
	"github.com/go-chi/chi/middleware"
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

type responseSegment struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

func AddSegmentHandler(log *slog.Logger, segmentAdder SegmentAdder) http.HandlerFunc {
	type request struct {
		Name        string `json:"name" validate:"required,min=4,max=255"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.AddSegmentHandler"

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

		if err := validator.New().Struct(req); err != nil {
			httpserver.RespondWithValidateError(w, log, err)
			return
		}

		if len(req.Description) > maxDescriptionLength {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Description is too long", log)
			return
		}

		req.Name = usecases_segments.FormatSegmnetName(req.Name)

		if _, err := segmentAdder.GetSegmentByName(r.Context(), req.Name); err == nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Segment with such name already exists", log)
			return
		}

		addedSegment, err := segmentAdder.AddSegment(r.Context(), models.AddSegmentParams{
			Name:        req.Name,
            Description: sql.NullString{
				String: req.Description,
				Valid: true,
			},
		})
		if err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not create segment", log)
			return
		}

		resp := responseSegment{
			Name:        addedSegment.Name,
			Description: addedSegment.Description.String,
			Created_At:  addedSegment.CreatedAt,
			Updated_At:  addedSegment.UpdatedAt,
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, resp)
	}
}

func DeleteSegmentHandler(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	type request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.v1.DeleteSegmentHandler"

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

		req.Name = usecases_segments.FormatSegmnetName(req.Name)

		if _, err := segmentDeleter.GetSegmentByName(r.Context(), req.Name); err != nil {
			httpserver.RespondWithError(w, http.StatusBadRequest, "Invalid segment name", log)
			return
		}

		if err = segmentDeleter.DeleteSegment(r.Context(), req.Name); err != nil {
			log.Error(err.Error())

			httpserver.RespondWithError(w, http.StatusInternalServerError, "Could not delete segment:", log)
			return
		}

		httpserver.RespondWithJSON(w, http.StatusOK, log, struct{}{})
	}
}