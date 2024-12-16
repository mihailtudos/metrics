package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mihailtudos/metrics/internal/domain/metrics"
	"github.com/mihailtudos/metrics/internal/infrastructure/http/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlePOSTMetric(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		method     string
		wantStatus int
		setupMock  func(*mocks.MetricsStore)
	}{
		{
			name:       "Success gauge metric",
			path:       "/update/gauge/Alloc/1.1",
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
			setupMock: func(m *mocks.MetricsStore) {
				m.On("Store", mock.MatchedBy(func(metric metrics.Metric) bool {
					return metric.ID == "Alloc" && metric.MType == metrics.GaugeType
				})).Return(nil)
			},
		},
		{
			name:       "Success counter metric",
			path:       "/update/counter/PollCount/1",
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
			setupMock: func(m *mocks.MetricsStore) {
				m.On("Store", mock.MatchedBy(func(metric metrics.Metric) bool {
					return metric.ID == "PollCount" && metric.MType == metrics.CounterType
				})).Return(nil)
			},
		},
		{
			name:       "Invalid metric type",
			path:       "/update/invalid/Test/1.1",
			method:     http.MethodPost,
			wantStatus: http.StatusBadRequest,
			setupMock:  func(m *mocks.MetricsStore) {},
		},
		{
			name:       "Invalid gauge value",
			path:       "/update/gauge/Alloc/notanumber",
			method:     http.MethodPost,
			wantStatus: http.StatusBadRequest,
			setupMock:  func(m *mocks.MetricsStore) {},
		},
		{
			name:       "Wrong method",
			path:       "/update/gauge/Alloc/1.1",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
			setupMock:  func(m *mocks.MetricsStore) {},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := mocks.NewMetricsStore(t)
			tc.setupMock(store)

			handler := NewHandler(store)
			r := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			handler.HandlePOSTMetric(w, r)
			res := w.Result()
			
			assert.Equal(t, tc.wantStatus, res.StatusCode)

			res.Body.Close()
		})
	}
}
