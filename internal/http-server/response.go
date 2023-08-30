package httpserver

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AlexZahvatkin/segments-users-service/internal/utils/validator"
	"github.com/go-playground/validator/v10"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, log *slog.Logger) {
	if code > 499 { 
		log.Error("Responding with 5XX error:", msg)
	}

	log.Error(msg)

	type errResponse struct {
		Error string `json:"error"`	
	}

	RespondWithJSON(w, code, log, errResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, log *slog.Logger, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Error("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func RespondWithCSV(w http.ResponseWriter, code int, log *slog.Logger, payload [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=report.csv")

	writer := csv.NewWriter(w)
	if err := writer.WriteAll(payload); err!= nil {
		log.Error("Failed to write CSV response: %v", payload)
	}
}

func RespondWithValidateError(w http.ResponseWriter, log *slog.Logger, err error) {
	validateErr := err.(validator.ValidationErrors)
	msg := validator_utils.ValidationError(validateErr)
	RespondWithError(w, http.StatusBadRequest, msg, log)
}