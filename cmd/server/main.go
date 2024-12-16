package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mihailtudos/metrics/internal/domain/metrics"
)

var ErrMissingMetricType = errors.New("missing metric type")
var ErrMissingMetricName = errors.New("missing metric name")
var ErrInvalidMetricValue = errors.New("invalid metric value")

type MemStorage struct {
	Metrics map[string]metrics.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]metrics.Metric),
	}
}

var Stroage *MemStorage

func main() {
	router := http.NewServeMux()
	Stroage = NewMemStorage()

	router.HandleFunc("POST /update/", handleMetrics)

	log.Println("Server started 🔥")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/update/")
	parts := strings.Split(url, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusNotFound)
		return
	}

	if parts[0] != "counter" && parts[0] != "gauge" {
		http.Error(w, ErrMissingMetricType.Error(), http.StatusBadRequest)
		return
	}

	if parts[1] == "" {
		http.Error(w, ErrMissingMetricName.Error(), http.StatusNotFound)
		return
	}

	metric := metrics.Metric{
		ID:    parts[1],
		MType: metrics.MetricType(parts[0]),
	}

	if parts[0] == "gauge" {
		val, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			http.Error(w, ErrInvalidMetricValue.Error(), http.StatusBadRequest)
			return
		}

		metric.Value = &val
	}

	if parts[0] == "counter" {
		val, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, ErrInvalidMetricValue.Error(), http.StatusBadRequest)
			return
		}

		metric.Delta = &val
	}

	Stroage.Metrics[metric.ID] = metric
}