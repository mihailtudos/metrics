package main

import (
	"log"
	"net/http"

	"github.com/mihailtudos/metrics/internal/infrastructure/http/handlers"
	"github.com/mihailtudos/metrics/internal/infrastructure/persistence/store"
)

func main() {
	router := http.NewServeMux()
	store := store.NewMetricStore()
	handlers := handlers.NewHandler(store)

	router.HandleFunc("POST /update/", handlers.HandlePOSTMetric)

	log.Println("Server started 🔥")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
