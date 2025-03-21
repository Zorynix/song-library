package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zorynix/song-library/config"
	_ "github.com/Zorynix/song-library/docs"
	logger "github.com/Zorynix/song-library/internal/logger"
	"github.com/Zorynix/song-library/internal/repo"
	v1 "github.com/Zorynix/song-library/internal/routes/http/v1"
	"github.com/Zorynix/song-library/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(configPath string) error {
	fmt.Println("Loading config from", configPath)
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to load config")
	}
	fmt.Println("Config loaded successfully")

	logger.SetupLogger(cfg.Log.Level)

	db, err := sqlx.Connect("postgres", cfg.PG.URL)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer db.Close()

	if err := RunMigrations(cfg.PG.URL); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	repos := repo.NewRepositories(db)
	services := services.NewServices(services.ServicesDependencies{
		Repos:       repos,
		MusicAPIURL: cfg.MusicAPI.URL,
	})
	handler := v1.NewHandler(services)

	r := chi.NewRouter()
	r.Route("/api/v1", handler.Register)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	apiServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Prometheus.MetricsPort),
		Handler: promhttp.Handler(),
	}

	go func() {
		logger.Logger.Info().Msgf("Starting API server on port %d", cfg.Server.Port)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal().Err(err).Msg("Failed to start API server")
		}
	}()

	go func() {
		logger.Logger.Info().Msgf("Starting metrics server on port %d", cfg.Prometheus.MetricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal().Err(err).Msg("Failed to start metrics server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info().Msg("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Logger.Fatal().Err(err).Msg("API server forced to shutdown")
	}
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Metrics server forced to shutdown")
	}

	logger.Logger.Info().Msg("Servers exited")

	return nil
}
