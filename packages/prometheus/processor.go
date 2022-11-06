package prometheus

import (
	"fmt"
	"strings"

	"github.com/fusakla/autograf/packages/model"
)

func guessMetricType(metric *model.Metric) {
	if strings.HasSuffix(metric.Name, "_total") {
		metric.MetricType = model.MetricTypeCounter
	}
	if strings.HasSuffix(metric.Name, "_info") || strings.HasSuffix(metric.Name, "_labels") {
		metric.MetricType = model.MetricTypeInfo
	}
	if metric.MetricType == model.MetricTypeHistogram && (strings.HasSuffix(string(metric.Name), "_sum") || strings.HasSuffix(string(metric.Name), "_count")) {
		metric.MetricType = model.MetricTypeCounter
	}
}

func guessMetricUnit(metric *model.Metric) {
	if strings.Contains(metric.Name, "_cpu_seconds") {
		metric.Unit = model.MetricUnitNone
	} else if strings.Contains(metric.Name, "_seconds") {
		metric.Unit = model.MetricUnitSeconds
	} else if strings.Contains(metric.Name, "_bytes") {
		metric.Unit = model.MetricUnitsBytes
	} else if strings.Contains(metric.Name, "_volt") {
		metric.Unit = model.MetricUnitsVolt
	} else if strings.Contains(metric.Name, "_ampere") {
		metric.Unit = model.MetricUnitsAmpere
	} else if strings.Contains(metric.Name, "_hertz") {
		metric.Unit = model.MetricUnitsHertz
	} else if strings.Contains(metric.Name, "_celsius") {
		metric.Unit = model.MetricUnitsCelsius
	}

}

func dropCreatedMetrics(metrics map[string]*model.Metric) {
	for k := range metrics {
		metricName := strings.TrimSuffix(k, "_created")
		if k == metricName {
			continue
		}
		if _, ok := metrics[metricName]; ok {
			delete(metrics, k)
			continue
		}
		if _, ok := metrics[metricName+"_total"]; ok {
			delete(metrics, k)
			continue
		}
		if _, ok := metrics[metricName+"_count"]; ok {
			delete(metrics, k)
			continue
		}
	}
}

func ProcessMetrics(metrics map[string]*model.Metric) error {
	var err error
	dropCreatedMetrics(metrics)
	for _, metric := range metrics {
		guessMetricUnit(metric)
		guessMetricType(metric)
		metric.Config, err = model.LoadConfigFromHelp(metric.Help)
		if err != nil {
			return fmt.Errorf("failed to parse autograf config in metric HELP: %w", err)
		}
	}
	return nil
}
