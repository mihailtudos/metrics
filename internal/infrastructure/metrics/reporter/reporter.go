package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	domain "github.com/mihailtudos/metrics/internal/domain/metrics"
)

type MetricsReporter struct {
	client    *http.Client
	ServerURL string
}

func NewMetricsReporter(ServerURL string) *MetricsReporter {
	return &MetricsReporter{
		client:    &http.Client{},
		ServerURL: ServerURL,
	}
}

func (m *MetricsReporter) ReportMetrics(metrics map[string]interface{}) {
	for key, val := range metrics {
		var metric domain.Metric

		intVal, ok := val.(int64)
		if ok {
			metric = domain.Metric{
				ID:    key,
				MType: "counter",
				Delta: &intVal,
			}
		} else {
			floatVal, ok := val.(float64)
			if ok {
				metric = domain.Metric{
					ID:    key,
					MType: "gauge",
					Value: &floatVal,
				}
			}
		}

		url := fmt.Sprintf("%s/update/", m.ServerURL)
		bData, err := json.Marshal(metric)
		if err != nil {
			log.Printf("failed to marshal metric: %v", err)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, url, io.Reader(bytes.NewBuffer(bData)))
		if err != nil {
			log.Printf("failed to create request: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := m.client.Do(req)
		if err != nil {
			log.Printf("failed to send request: %v", err)
			continue
		}

		log.Println("reported: ", url, "; status ", res.StatusCode)
		res.Body.Close()
	}
}
