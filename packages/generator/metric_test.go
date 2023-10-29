package generator

import (
	"testing"

	"gotest.tools/assert"
)

func TestMetric_PromQlQuery(t *testing.T) {

	type tc struct {
		name          string
		metric        Metric
		selector      string
		rangeSelector string
		expectedQuery string
	}
	cases := []tc{
		{
			name: "gauge",
			metric: Metric{
				MetricType: MetricTypeGauge,
				Name:       "queue_size",
				Config:     MetricConfig{},
			},
			selector:      `{foo="bar"}`,
			rangeSelector: "5m",
			expectedQuery: `queue_size{foo="bar"}`,
		},
		{
			name: "counter",
			metric: Metric{
				MetricType: MetricTypeCounter,
				Name:       "counter_total",
				Config:     MetricConfig{},
			},
			selector:      `{foo="bar"}`,
			rangeSelector: "5m",
			expectedQuery: `rate(counter_total{foo="bar"}[5m])`,
		},
		{
			name: "histogram",
			metric: Metric{
				MetricType: MetricTypeHistogram,
				Name:       "histogram_count",
				Config:     MetricConfig{},
			},
			selector:      `{foo="bar"}`,
			rangeSelector: "5m",
			expectedQuery: `sum(rate(histogram_count{foo="bar"}[5m])) by (le)`,
		},
		{
			name: "summary",
			metric: Metric{
				MetricType: MetricTypeSummary,
				Name:       "summary",
				Config:     MetricConfig{},
			},
			selector:      `{foo="bar"}`,
			rangeSelector: "5m",
			expectedQuery: `sum(summary{foo="bar"}) by (quantile)`,
		},
		{
			name: "info",
			metric: Metric{
				MetricType: MetricTypeInfo,
				Name:       "app_info",
				Config:     MetricConfig{},
			},
			selector:      `{foo="bar"}`,
			rangeSelector: "5m",
			expectedQuery: `app_info{foo="bar"}`,
		},
		{
			name: "gauge with time",
			metric: Metric{
				MetricType: MetricTypeGauge,
				Name:       "last_timestamp_seconds",
				Config:     MetricConfig{},
			},
			selector:      `{foo="bar"}`,
			rangeSelector: "5m",
			expectedQuery: `time() - last_timestamp_seconds{foo="bar"}`,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedQuery, tt.metric.PromQlQuery(tt.selector, tt.rangeSelector))
		})
	}
}
