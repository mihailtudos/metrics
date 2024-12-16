package reporter

import (
	"fmt"
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
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	// http://localhost:8080
	// Content-Type: text/plain

	for key, val := range metrics {
		t := domain.GaugeType

		_, ok := val.(int64)
		if ok {
			t = domain.CounterType
		}

		url := fmt.Sprintf("%s/update/%s/%s/%v", m.ServerURL, t, key, val)

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Printf("failed to create request: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		res, err := m.client.Do(req)
		if err != nil {
			log.Printf("failed to send request: %v", err)
			continue
		}

		log.Println("reported: ", url, "; status ", res.StatusCode)
		res.Body.Close()
	}
}
