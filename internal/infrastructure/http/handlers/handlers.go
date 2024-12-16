package handlers

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

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=MetricsStore
type MetricsStore interface {
	Store(metric metrics.Metric) error
}

type Handler struct {
	MetricsStore
}

func NewHandler(Store MetricsStore) *Handler {
	return &Handler{
		MetricsStore: Store,
	}
}

func (h *Handler) HandlePOSTMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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

	if err := h.Store(metric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to store metric")
		return
	}

	w.WriteHeader(http.StatusOK)
}
