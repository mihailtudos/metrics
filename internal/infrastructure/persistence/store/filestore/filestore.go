package filestore

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/mihailtudos/metrics/internal/domain/metrics"
)

type FileStorage struct {
	Metrics       metrics.Metrics
	Encoder       *json.Encoder
	Decoder       *json.Decoder
	StoreFile     *os.File
	Logger        *slog.Logger
	filePath      string
	restore       bool
	storeInterval time.Duration
	ticker        *time.Ticker
	cancel        context.Context
}

func NewFileStore(ctx context.Context, logger *slog.Logger, storeInterval time.Duration, fileStorePath string, restore bool) *FileStorage {
	file, err := os.OpenFile(fileStorePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.ErrorContext(context.Background(), "failed to open file", slog.String("error", err.Error()))
		panic(err)
	}

	store := &FileStorage{
		Logger:        logger,
		Encoder:       json.NewEncoder(file),
		Decoder:       json.NewDecoder(file),
		StoreFile:     file,
		Metrics:       map[string]metrics.Metric{},
		filePath:      fileStorePath,
		restore:       restore,
		storeInterval: storeInterval,
		cancel:        ctx,
	}

	if storeInterval != 0 {
		store.ticker = time.NewTicker(storeInterval)
		go store.startTicker(ctx)
	}

	if store.restore {
		store.LoadInFile()
	}

	return store
}

func (m *FileStorage) Store(metric metrics.Metric) error {
	if metric.MType == metrics.CounterType {
		_, ok := m.Metrics[metric.ID]
		if ok {
			*m.Metrics[metric.ID].Delta += *metric.Delta
		} else {
			m.Metrics[metric.ID] = metric
		}

		return nil
	}

	m.Metrics[metric.ID] = metric

	m.SaveFile(m.storeInterval != 0)

	return nil
}

func (m *FileStorage) GetAllMetrics() ([]metrics.Metric, error) {
	var metrics []metrics.Metric
	for _, metric := range m.Metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *FileStorage) GetOneMetric(metricName string) (metrics.Metric, error) {
	metric, ok := m.Metrics[metricName]
	if !ok {
		return metrics.Metric{}, metrics.ErrMetricNotFound
	}

	return metric, nil
}

func (m *FileStorage) LoadInFile() error {
	m.Logger.InfoContext(context.Background(), "loading metrics from file")
	return m.Decoder.Decode(&m.Metrics)
}

func (m *FileStorage) Close() error {
	m.Logger.InfoContext(context.Background(), "closing file")

	m.ticker.Stop()
	return m.StoreFile.Close()
}

func (m *FileStorage) SaveFile(shouldSync bool) error {
	if !shouldSync {
		m.Logger.InfoContext(context.Background(), "skiping file sync")
		return nil
	}

	m.Logger.InfoContext(context.Background(), "saving metrics to file")
	return m.Encoder.Encode(m.Metrics)
}

func (m *FileStorage) startTicker(ctx context.Context) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-m.ticker.C:
				m.SaveFile(true)
			case <-ctx.Done():
				m.Logger.InfoContext(ctx, "stopping ticker")
				m.SaveFile(true)
				m.Close()
				return
			}
		}
	}()

	<-done
}
