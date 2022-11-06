package model

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

type MetricUnit string

const (
	MetricUnitNone     MetricUnit = "none"
	MetricUnitSeconds  MetricUnit = "s"
	MetricUnitsBytes   MetricUnit = "decbytes"
	MetricUnitsAmpere  MetricUnit = "amp"
	MetricUnitsVolt    MetricUnit = "volt"
	MetricUnitsHertz   MetricUnit = "rothz"
	MetricUnitsCelsius MetricUnit = "celsius"
)

type Metric struct {
	MetricType MetricType
	Name       string
	Help       string
	Unit       MetricUnit
	Comment    string
	Config     PanelConfig
}
