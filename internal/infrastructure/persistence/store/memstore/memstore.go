package memstore

import "github.com/mihailtudos/metrics/internal/domain/metrics"

type MemStorage struct {
	Metrics map[string]metrics.Metric
}

func NewMemStore() *MemStorage {
	return &MemStorage{
		Metrics: map[string]metrics.Metric{},
	}
}

func (m *MemStorage) Store(metric metrics.Metric) error {
	m.Metrics[metric.ID] = metric

	return nil
}
