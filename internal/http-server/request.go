package httpserver

// func getUserIdFromParams(w http.ResponseWriter, r *http.Request, log *slog.Logger) (int64, error) {
// 	s := chi.URLParam(r, "userId")
// 	userId, err := strconv.ParseInt(s, 10, 64)
// 	if err != nil {
// 		httpserver.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("User id must be a number: %v", err), log)
// 		return 0, err
// 	}
// 	return userId, nil
// }

// func getTimeFromParams(w http.ResponseWriter, r *http.Request, log *slog.Logger, paramName string) (time.Time, error) {
// 	s := chi.URLParam(r, "userId")
// 	time, err := strconv.ParseInt(s, 10, 64)
	
// }