package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/AlexZahvatkin/segments-users-service/config"
	"github.com/AlexZahvatkin/segments-users-service/internal/storage/database"
	v1 "github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1"
	"github.com/AlexZahvatkin/segments-users-service/internal/lib/logger"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting segments-users-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
 
	log.Info("Initializing postgres...")
	queries := initDb(getDbURL(cfg), log)
	log.Debug(getDbURL(cfg))

	log.Info("Initializing routers...")
	router := v1.InitRouters(queries, log, cfg)

	srv := &http.Server{
		Addr: cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {  
		log.Error("failed to start server", err)
	}

	log.Error("server stopped")
}

func initDb(dbURL string, log *slog.Logger) *database.Queries {
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
        log.Error("can not connect to a database", sl.Err(err))
		os.Exit(1)
    }
	queries := database.New(conn)
	return queries
}

func getDbURL(cfg *config.Config) string{
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Database.User, cfg.Database.Password, 
	cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}