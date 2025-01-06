package store

import (
	"context"
	"log/slog"
	"time"

	"github.com/mihailtudos/metrics/internal/domain/metrics"
	"github.com/mihailtudos/metrics/internal/infrastructure/persistence/store/filestore"
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

func NewMetricStore(ctx context.Context, logger *slog.Logger, storeInterval time.Duration, fileStorePath string, restore bool) *MetricStore {
	var store MetricsStorer
	if storeInterval >= 0 {
		logger.InfoContext(context.Background(), "using file store")
		store = filestore.NewFileStore(ctx, logger, storeInterval, fileStorePath, restore)
	}

	if store == nil {
		store = memstore.NewMemStore(logger)
	}

	return &MetricStore{
		MetricsStorer: store,
	}
}
