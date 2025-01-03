package store

import (
	"github.com/mihailtudos/metrics/internal/domain/metrics"
	"github.com/mihailtudos/metrics/internal/infrastructure/persistence/store/memstore"
)

type MetricsStorer interface {
	Store(metrics.Metric) error
	GetAllMetrics() ([]metrics.Metric, error)
	GetOneMetric(metricName string) (metrics.Metric, error)
}

type MetricStore struct {
	MetricsStorer
}

func NewMetricStore() *MetricStore {
	store := memstore.NewMemStore()
	return &MetricStore{
		MetricsStorer: store,
	}
}
