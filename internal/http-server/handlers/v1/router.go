package v1

import (
	"log/slog"

	"github.com/AlexZahvatkin/segments-users-service/internal/database"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/segments"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/users"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/middleware/mwlogger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func InitRouters(queries *database.Queries, log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*, http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))
	
	v1Router := chi.NewRouter()

	v1Router.Use(middleware.RequestID)
	v1Router.Use(middleware.Logger)
	v1Router.Use(mwlogger.New(log))
	v1Router.Use(middleware.Recoverer)
	v1Router.Use(middleware.URLFormat)

	v1Router.Post("/users", users.AddUserHandler(log, queries))
	v1Router.Delete("/users", users.DeleteUserHandler(log, queries))
	v1Router.Post("/segments", segments.AddSegmemtHandler(log, queries))
	v1Router.Delete("/segments", segments.DeleteSegmemtHandler(log, queries))

	router.Mount("/v1", v1Router)

	return router
}