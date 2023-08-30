package handlers

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func SetLogger(log *slog.Logger, ctx context.Context, op string) {
	log = log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)
}