package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func GetUserIdFromParams(w http.ResponseWriter, r *http.Request, log *slog.Logger) (int64, error) {
	s := chi.URLParam(r, "userId")
	log.Debug("Param is " + s)
	if s == "" {
		RespondWithError(w, http.StatusBadRequest, "You must provide userId", log)
		return -1, errors.New("No parameter provided")
	}
	userId, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("User id must be a number: %v", err), log)
		return -1, err
	}
	return userId, nil
}

func GetTimeFromParams(w http.ResponseWriter, r *http.Request, log *slog.Logger, paramName string, timeFormat string) (time.Time, error) {
	s := r.URL.Query().Get(paramName)
	if s == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("You must provide %s", paramName), log)
		return time.Time{}, errors.New("No parameter provided")
	}
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Wrong datetime format: %v", err), log)
		return time.Time{}, err
	}
	return t, nil
}

func DecodeRequsetBody[t any](w http.ResponseWriter, r *http.Request, req t, log *slog.Logger) (t, error) {
	decoder := json.NewDecoder(r.Body) 
	err := decoder.Decode(&req)
	if err!= nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err), log)
		return req, err
	}
	return req, nil
}