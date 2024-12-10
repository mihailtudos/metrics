package main

import (
	"log"
	"net/http"
	"strings"
)

type Metric struct {
	ID    string
	MType string
	Delta int64
	Value float64
}

type MemStorage struct {
	Metrics map[string]Metric
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/update/", handleMetrics)

	log.Println("Server started 🔥")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	// /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	url := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := strings.Split(url, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	
}
