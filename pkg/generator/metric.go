package generator

import "regexp"

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
	Config     MetricConfig
	Threshold  *Metric
}

var timeMetricsRegexp = regexp.MustCompile(`.+(_time|_time_seconds|_timestamp|_timestamp_seconds|_update|_started|_last_seen)$`)

func (m *Metric) PromQlQuery(selector string, rangeSelector string) string {
	query := ""
	switch m.MetricType {
	case MetricTypeGauge:
		query = metricWithSelector(m.Name, selector)
		if timeMetricsRegexp.MatchString(m.Name) {
			query = "time() - " + query
			m.Unit = MetricUnitSeconds
		}
	case MetricTypeCounter:
		if timeMetricsRegexp.MatchString(m.Name) {
			query = "time() - " + m.Name
			m.Unit = MetricUnitSeconds
		} else {
			query = rateCounterQuery(metricWithSelector(m.Name, selector), rangeSelector)
		}
	case MetricTypeHistogram:
		query = rateCounterQuery(metricWithSelector(m.Name, selector), rangeSelector)
		m.Config.AggregateBy = append(m.Config.AggregateBy, "le")
		if m.Config.Aggregation == "" {
			m.Config.Aggregation = "sum"
		}
	case MetricTypeSummary:
		query = metricWithSelector(m.Name, selector)
		m.Config.AggregateBy = append(m.Config.AggregateBy, "quantile")
		if m.Config.Aggregation == "" {
			m.Config.Aggregation = "sum"
		}
	default:
		query = metricWithSelector(m.Name, selector)
	}
	if m.Config.Aggregation != "" {
		query = aggregateQuery(query, m.Config.Aggregation, m.Config.AggregateBy)
	}
	return query
}
