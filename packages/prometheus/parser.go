package prometheus

import (
	"io"

	"github.com/prometheus/prometheus/model/textparse"
)

func ParseMetricsText(text []byte, openMetrics bool) (map[string]Metric, error) {
	var (
		metrics       = make(map[string]Metric)
		currentMetric Metric
	)
	contentType := "text"
	if openMetrics {
		contentType = "application/openmetrics-text"
	}
	p, err := textparse.New(text, contentType)
	if err != nil {
		return nil, err
	}
	storeCurrentMetric := func(name string) {
		metrics[currentMetric.Name] = currentMetric
		currentMetric = Metric{Name: name, Unit: guessMetricUnit(name)}
	}

	for {
		var entryType textparse.Entry
		if entryType, err = p.Next(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch entryType {
		case textparse.EntryType:
			mName, mType := p.Type()
			if currentMetric.Name != string(mName) {
				storeCurrentMetric(string(mName))
			}
			currentMetric.Name = string(mName)
			currentMetric.MetricType = MetricType(mType)
		case textparse.EntryHelp:
			mName, mHelp := p.Help()
			if currentMetric.Name != string(mName) {
				storeCurrentMetric(string(mName))
			}
			currentMetric.Name = string(mName)
			currentMetric.Help = string(mHelp)
		case textparse.EntryUnit:
			_, mUnit := p.Unit()
			currentMetric.Unit = MetricUnit(mUnit)
		default:
		}
	}
	metrics[currentMetric.Name] = currentMetric
	delete(metrics, "")
	return metrics, nil
}
