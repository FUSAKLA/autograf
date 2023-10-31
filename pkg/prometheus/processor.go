package prometheus

import (
	"fmt"
	"strings"

	"github.com/fusakla/autograf/pkg/generator"
)

func guessMetricType(metric *generator.Metric) {
	if strings.HasSuffix(metric.Name, "_total") {
		metric.MetricType = generator.MetricTypeCounter
	}
	if strings.HasSuffix(metric.Name, "_info") || strings.HasSuffix(metric.Name, "_labels") {
		metric.MetricType = generator.MetricTypeInfo
	}
	if metric.MetricType == generator.MetricTypeHistogram && (strings.HasSuffix(string(metric.Name), "_sum") || strings.HasSuffix(string(metric.Name), "_count")) {
		metric.MetricType = generator.MetricTypeCounter
	}
}

func guessMetricUnit(metric *generator.Metric) {
	if strings.Contains(metric.Name, "_cpu_seconds") {
		metric.Unit = generator.MetricUnitNone
	} else if strings.Contains(metric.Name, "_seconds") {
		metric.Unit = generator.MetricUnitSeconds
	} else if strings.Contains(metric.Name, "_bytes") {
		metric.Unit = generator.MetricUnitsBytes
	} else if strings.Contains(metric.Name, "_volt") {
		metric.Unit = generator.MetricUnitsVolt
	} else if strings.Contains(metric.Name, "_ampere") {
		metric.Unit = generator.MetricUnitsAmpere
	} else if strings.Contains(metric.Name, "_hertz") {
		metric.Unit = generator.MetricUnitsHertz
	} else if strings.Contains(metric.Name, "_celsius") {
		metric.Unit = generator.MetricUnitsCelsius
	}

}

func dropCreatedMetrics(metrics map[string]*generator.Metric) {
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

func ProcessMetrics(metrics map[string]*generator.Metric) error {
	var err error
	dropCreatedMetrics(metrics)
	for _, metric := range metrics {
		guessMetricUnit(metric)
		guessMetricType(metric)
		metric.Config, err = generator.LoadConfigFromHelp(metric.Help)
		if err != nil {
			return fmt.Errorf("failed to parse autograf config in metric HELP: %w", err)
		}
		metric.Threshold = metric.Config.ThresholdMetric()
	}
	return nil
}
