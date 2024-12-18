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

func (m *MemStorage) GetAllMetrics() ([]metrics.Metric, error) {
	var metrics []metrics.Metric
	for _, metric := range m.Metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *MemStorage) GetOneMetric(metricName string) (metrics.Metric, error) {
	metric, ok := m.Metrics[metricName]
	if !ok {
		return metrics.Metric{}, metrics.ErrMetricNotFound
	}

	return metric, nil
}
