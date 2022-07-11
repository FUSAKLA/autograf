package prometheus

import (
	"strings"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func guessMetricType(metricName string) MetricType {
	if strings.HasSuffix(metricName, "_total") {
		return MetricTypeCounter
	}
	if strings.HasSuffix(metricName, "_info") {
		return MetricTypeInfo
	}
	return MetricTypeGauge
}

func guessMetricUnit(metricName string) string {
	if strings.Contains(metricName, "_seconds") {
		return "s"
	}
	if strings.Contains(metricName, "_bytes") {
		return "decbytes"
	}
	return "none"
}

func guessMetricMetadata(metricName string) []v1.Metadata {
	return []v1.Metadata{
		{
			Type: v1.MetricType(guessMetricType(metricName)),
			Help: "",
			Unit: guessMetricUnit(metricName),
		},
	}
}
