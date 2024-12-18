package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mihailtudos/metrics/internal/infrastructure/http/handlers"
	"github.com/mihailtudos/metrics/internal/infrastructure/persistence/store"
)

func main() {
	router := chi.NewRouter()
	store := store.NewMetricStore()
	handlers := handlers.NewHandler(store)

	router.Post("/update/{type}/{name}/{value}", handlers.HandlePOSTMetric)
	router.Get("/", handlers.HandleShowAllMetrics)
	router.Get("/value/{type}/{name}", handlers.HandleShowMetricValue)

	log.Println("Server started 🔥")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
