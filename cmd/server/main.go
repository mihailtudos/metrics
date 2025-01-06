package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/mihailtudos/metrics/internal/infrastructure/config/logger"
	"github.com/mihailtudos/metrics/internal/infrastructure/config/server"
	"github.com/mihailtudos/metrics/internal/infrastructure/http/handlers"
	"github.com/mihailtudos/metrics/internal/infrastructure/http/middleware"
	"github.com/mihailtudos/metrics/internal/infrastructure/persistence/store"
	utils "github.com/mihailtudos/metrics/utils/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srvConfig := server.NewServerConfig()
	logger := logger.NewLogger()
	logger.DebugContext(ctx, "server configuration", slog.Any("config", srvConfig))

	router := chi.NewRouter()
	store := store.NewMetricStore(
		ctx,
		logger,
		srvConfig.StoreInterval,
		srvConfig.FileStorePath,
		srvConfig.Restore,
	)
	handlers := handlers.NewHandler(store, logger)

	router.Use(middleware.WithLogger(logger))
	router.Use(middleware.WithCompress)

	router.Post("/update/{type}/{name}/{value}", handlers.HandlePOSTMetric)
	router.Get("/", handlers.HandleShowAllMetrics)
	router.Get("/value/{type}/{name}", handlers.HandleShowMetricValue)

	router.Post("/update/", handlers.HandlePOSTMetricWithJSON)
	router.Post("/value/", handlers.HandleShowMetricValueWithJSON)

	logger.InfoContext(ctx, "server started at", slog.String("address", srvConfig.Address))

	srv := &http.Server{
		Addr:    srvConfig.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(context.Background(), "failed to start server", utils.ErrValue(err))

			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	logger.InfoContext(ctx, "got interruption signal")

	if err := srv.Shutdown(context.TODO()); err != nil {
		logger.InfoContext(ctx, "server shutdown returned an err")
	}

	logger.InfoContext(ctx, "exting")
}
