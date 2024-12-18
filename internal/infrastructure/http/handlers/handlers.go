package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/mihailtudos/metrics/internal/domain/metrics"
)

var ErrMissingMetricType = errors.New("missing metric type")
var ErrMissingMetricName = errors.New("missing metric name")
var ErrInvalidMetricValue = errors.New("invalid metric value")

//go:embed templates/metrics.index.html
var fs embed.FS

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=MetricsStore
type MetricsStore interface {
	Store(metric metrics.Metric) error
	GetAllMetrics() ([]metrics.Metric, error)
	GetOneMetric(metricName string) (metrics.Metric, error)
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
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

	if mType == "" || mType != string(metrics.CounterType) && mType != string(metrics.GaugeType) {
		http.Error(w, ErrMissingMetricType.Error(), http.StatusBadRequest)
		return
	}

	if mName == "" {
		http.Error(w, ErrMissingMetricName.Error(), http.StatusNotFound)
		return
	}

	metric := metrics.Metric{
		ID:    mName,
		MType: metrics.MetricType(mType),
	}

	if mType == string(metrics.GaugeType) {
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, ErrInvalidMetricValue.Error(), http.StatusBadRequest)
			return
		}

		metric.Value = &val
	}

	if mType == string(metrics.CounterType) {
		val, err := strconv.ParseInt(mValue, 10, 64)
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

func (h *Handler) HandleShowAllMetrics(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(fs, "templates/metrics.index.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	metrics, err := h.MetricsStore.GetAllMetrics()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	tmpl.Execute(w, metrics)
}

func (h *Handler) HandleShowMetricValue(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")

	if metricType == "" ||
		(metricType != string(metrics.CounterType) && metricType != string(metrics.GaugeType)) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	metric, err := h.MetricsStore.GetOneMetric(name)
	if err != nil {
		if errors.Is(err, metrics.ErrMetricNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res := ""
	if metric.Delta != nil {
		res = fmt.Sprintf("%v", *metric.Delta)
	}

	if metric.Value != nil {
		res = fmt.Sprintf("%v", *metric.Value)
	}

	w.WriteHeader(200)
	w.Write([]byte(res))
}
