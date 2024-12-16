package metrics

type MetricType string

const (
    CounterType MetricType = "counter"
    GaugeType   MetricType = "gauge"
)

type Metric struct {
    ID    string     `json:"id"`
    MType MetricType `json:"type"`
    Delta *int64     `json:"delta,omitempty"`
    Value *float64   `json:"value,omitempty"`
}

// Value object for validation
type MetricValidation struct {
    IsValid bool
    Errors  []string
}

func (m *Metric) Validate() MetricValidation {
    validation := MetricValidation{IsValid: true}
    
    if m.ID == "" {
        validation.IsValid = false
        validation.Errors = append(validation.Errors, "metric ID cannot be empty")
    }

    switch m.MType {
    case CounterType:
        if m.Delta == nil {
            validation.IsValid = false
            validation.Errors = append(validation.Errors, "counter metric must have delta value")
        }
    case GaugeType:
        if m.Value == nil {
            validation.IsValid = false
            validation.Errors = append(validation.Errors, "gauge metric must have value")
        }
    default:
        validation.IsValid = false
        validation.Errors = append(validation.Errors, "invalid metric type")
    }

    return validation
}