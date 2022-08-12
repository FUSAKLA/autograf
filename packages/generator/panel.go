package generator

import (
	"fmt"

	"github.com/K-Phoen/grabana/heatmap"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/table"
	grabana_prometheus "github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/fusakla/autograf/packages/prometheus"
)

func newTimeSeriesPanel(dataSource, selector string, metric prometheus.Metric) row.Option {
	query := fmt.Sprintf("%s%s", metric.Name, selector)
	if metric.MetricType == "counter" {
		query = fmt.Sprintf("rate(%s%s[3m])", metric.Name, selector)
	}
	return row.WithTimeSeries(
		metric.Name,
		timeseries.Description(metric.Help),
		timeseries.Axis(
			axis.Unit(metric.Unit),
		),
		timeseries.DataSource(dataSource),
		timeseries.WithPrometheusTarget(
			query,
			grabana_prometheus.Format(grabana_prometheus.FormatTimeSeries),
			grabana_prometheus.Ref(metric.Name),
		),
	)
}

func newHeatmapPanel(dataSource, selector string, metric prometheus.Metric) row.Option {
	query := fmt.Sprintf("sum(rate(%s%s[3m])) by (le)", metric.Name, selector)
	return row.WithHeatmap(
		metric.Name,
		heatmap.Description(metric.Help),
		heatmap.HideZeroBuckets(),
		heatmap.DataSource(dataSource),
		heatmap.WithPrometheusTarget(
			query,
			grabana_prometheus.Format(grabana_prometheus.FormatHeatmap),
			grabana_prometheus.Ref(metric.Name),
		),
	)
}

func newInfoPanel(dataSource, selector string, metric prometheus.Metric) row.Option {
	return row.WithTable(
		metric.Name,
		table.Description(metric.Help),
		table.WithPrometheusTarget(
			fmt.Sprintf("%s%s", metric.Name, selector),
			grabana_prometheus.Format(grabana_prometheus.FormatTable),
			grabana_prometheus.Instant(),
			grabana_prometheus.Ref(metric.Name),
		),
	)
}

func newPanel(dataSource string, selector string, metric prometheus.Metric) row.Option {
	switch metric.MetricType {
	case "gauge":
		return newTimeSeriesPanel(dataSource, selector, metric)
	case "counter":
		return newTimeSeriesPanel(dataSource, selector, metric)
	case "summary":
		return newTimeSeriesPanel(dataSource, selector, metric)
	case "histogram":
		return newHeatmapPanel(dataSource, selector, metric)
	case "info":
		return newInfoPanel(dataSource, selector, metric)
	}
	metric.Comment = fmt.Sprintf("WARNING: Unknown metric type %s!\n%s", metric.MetricType, metric.Comment)
	return newTimeSeriesPanel(dataSource, selector, metric)
}
