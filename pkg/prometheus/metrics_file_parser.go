package prometheus

import (
	"io"
	"strings"

	"github.com/fusakla/autograf/pkg/generator"
	"github.com/prometheus/prometheus/model/textparse"
)

func newFileMetrics() fileMetrics {
	return fileMetrics{metrics: map[string]*generator.Metric{}}
}

type fileMetrics struct {
	metrics map[string]*generator.Metric
}

func (f *fileMetrics) add(metric generator.Metric) {
	m, ok := f.metrics[metric.Name]
	if !ok {
		f.metrics[metric.Name] = &metric
		return
	}
	if metric.Name != "" {
		m.Name = metric.Name
	}
	if metric.Help != "" {
		m.Help = metric.Help
	}
	if metric.Unit != "" {
		m.Unit = metric.Unit
	}
	if metric.MetricType != "" {
		m.MetricType = metric.MetricType
	}
}

func (f *fileMetrics) addHistograms(histograms []string) {
	for _, h := range histograms {
		m := f.metrics[h]
		f.add(generator.Metric{Name: m.Name + "_bucket", MetricType: generator.MetricTypeHistogram, Help: m.Help, Unit: m.Unit})
		f.add(generator.Metric{Name: m.Name + "_sum", MetricType: generator.MetricTypeHistogram, Help: m.Help, Unit: m.Unit})
		f.add(generator.Metric{Name: m.Name + "_count", MetricType: generator.MetricTypeHistogram, Help: m.Help, Unit: m.Unit})
		delete(f.metrics, h)
	}
}

func (f *fileMetrics) finalMetrics() map[string]*generator.Metric {
	delete(f.metrics, "")
	return f.metrics
}

func ParseMetricsText(text []byte, openMetrics bool) (map[string]*generator.Metric, error) {
	metrics := newFileMetrics()
	histograms := []string{}
	contentType := "text"
	if openMetrics {
		contentType = "application/openmetrics-text"
	}
	p, err := textparse.New(text, contentType, true)
	if err != nil {
		return nil, err
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
			if mType == textparse.MetricTypeHistogram {
				histograms = append(histograms, string(mName))
			}
			metrics.add(generator.Metric{Name: string(mName), MetricType: generator.MetricType(mType)})
		case textparse.EntryHelp:
			mName, mHelp := p.Help()
			metrics.add(generator.Metric{Name: string(mName), Help: string(mHelp)})
		case textparse.EntryUnit:
			mName, mUnit := p.Unit()
			metrics.add(generator.Metric{Name: string(mName), Unit: generator.MetricUnit(mUnit)})
		case textparse.EntrySeries:
			m, _, _ := p.Series()
			metrics.add(generator.Metric{Name: strings.Split(string(m), "{")[0]})
		default:
		}
	}
	metrics.addHistograms(histograms)
	return metrics.finalMetrics(), nil
}
