package prometheus

// MetricType models the type of a metric.
type MetricType string

const (
	// Possible values for MetricType
	MetricTypeCounter        MetricType = "counter"
	MetricTypeGauge          MetricType = "gauge"
	MetricTypeHistogram      MetricType = "histogram"
	MetricTypeGaugeHistogram MetricType = "gaugehistogram"
	MetricTypeSummary        MetricType = "summary"
	MetricTypeInfo           MetricType = "info"
	MetricTypeStateset       MetricType = "stateset"
	MetricTypeUnknown        MetricType = "unknown"
)

type Metric struct {
	MetricType MetricType
	Name       string
	Help       string
	Unit       string
	Comment    string
}

func (m Metric) String() string {
	return m.Name
}
