package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mihailtudos/metrics/internal/infrastructure/config/logger"
	"github.com/mihailtudos/metrics/internal/infrastructure/config/server"
	"github.com/mihailtudos/metrics/internal/infrastructure/http/handlers"
	"github.com/mihailtudos/metrics/internal/infrastructure/http/middleware"
	"github.com/mihailtudos/metrics/internal/infrastructure/persistence/store"
)

func main() {
	srvConfig := server.NewServerConfig()
	logger := logger.NewLogger()

	router := chi.NewRouter()
	store := store.NewMetricStore()
	handlers := handlers.NewHandler(store)

	router.Use(middleware.WithLogger(logger))

	router.Post("/update/{type}/{name}/{value}", handlers.HandlePOSTMetric)
	router.Get("/", handlers.HandleShowAllMetrics)
	router.Get("/value/{type}/{name}", handlers.HandleShowMetricValue)

	logger.InfoContext(nil, "Server started at ", slog.String("address", srvConfig.Address))
	if err := http.ListenAndServe(srvConfig.Address, router); err != nil {
		panic(err)
	}
}
