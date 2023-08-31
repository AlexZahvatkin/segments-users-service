package v1

import (
	"fmt"
	"log/slog"

	"github.com/AlexZahvatkin/segments-users-service/config"
	_ "github.com/AlexZahvatkin/segments-users-service/docs"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/segments"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/users"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/users_in_segments"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/middleware/mwlogger"
	"github.com/AlexZahvatkin/segments-users-service/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func InitRouters(storage storage.Storage, log *slog.Logger, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*, http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Use(middleware.RequestID)
	v1Router.Use(middleware.Logger)
	v1Router.Use(mwlogger.New(log))
	v1Router.Use(middleware.Recoverer)
	v1Router.Use(middleware.URLFormat)

	v1Router.Post("/segments/assign/{userId}", users_in_segments.SegmentsAssignHandler(log, storage))
	v1Router.Get("/segments/history/{userId}", users_in_segments.GetSegmentsHistoryByUser(log, storage))
	v1Router.Post("/segments/ttl/{userId}", users_in_segments.SegmentsAssignWithTTLInHoursHandler(log, storage))
	v1Router.Get("/segments/{userId}", users_in_segments.GetSegmentsForUserHandler(log, storage))
	v1Router.Post("/users", users.AddUserHandler(log, storage))
	v1Router.Delete("/users/{userId}", users.DeleteUserHandler(log, storage))
	v1Router.Post("/segments", segments.AddSegmentHandler(log, storage))
	v1Router.Delete("/segments", segments.DeleteSegmentHandler(log, storage))

	router.Mount("/v1", v1Router)

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port))))

	return router
}
