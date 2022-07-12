package generator

import (
	"fmt"
	"strings"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/fusakla/autograf/packages/prometheus"
)

func newTree(prefix string, level int) *metricsTree {
	return &metricsTree{
		prefix: prefix,
		level:  level,
		leafs:  make(map[string]*metricsTree),
	}
}

type metricsTree struct {
	leafs  map[string]*metricsTree
	prefix string
	metric *prometheus.Metric
	level  int
}

func (t *metricsTree) metrics() []*prometheus.Metric {
	metrics := []*prometheus.Metric{}
	if t.metric != nil {
		metrics = append(metrics, t.metric)
	}
	for _, l := range t.leafs {
		if l.metric != nil {
			metrics = append(metrics, l.metric)
		} else {
			metrics = append(metrics, l.metrics()...)
		}
	}
	return metrics
}

func (t *metricsTree) metricGroups() map[string][]*prometheus.Metric {
	others := []*prometheus.Metric{}
	groups := map[string][]*prometheus.Metric{}
	for _, l := range t.leafs {
		subtreeMetrics := l.metrics()
		if len(subtreeMetrics) < 3 && t.prefix != "" {
			others = append(others, subtreeMetrics...)
		} else {
			for k, v := range l.metricGroups() {
				groups[k] = v
			}
		}
	}
	if len(others) > 0 {
		groups[t.prefix] = others
	}
	return groups
}

func (t metricsTree) String() string {
	indent := strings.Repeat(" ", t.level)
	out := ""
	out += fmt.Sprintf(indent + t.prefix + "\n")
	if t.metric != nil {
		out += fmt.Sprintf(indent + "- " + t.metric.String())
	}
	for _, l := range t.leafs {
		out += "\n" + l.String()
	}
	return out
}

func (t *metricsTree) add(metric prometheus.Metric) {
	strippedMetricName := strings.TrimPrefix(strings.TrimPrefix(metric.Name, t.prefix), "_")
	if strippedMetricName == "" {
		t.metric = &metric
		return
	}
	parts := strings.SplitN(strippedMetricName, "_", 2)
	prefix := parts[0]
	if _, ok := t.leafs[prefix]; !ok {
		t.leafs[prefix] = newTree(strings.TrimPrefix(t.prefix+"_"+prefix, "_"), t.level+1)
	}
	t.leafs[prefix].add(metric)
}

func groupMetrics(metrics map[string]prometheus.Metric) map[string][]*prometheus.Metric {
	tree := newTree("", 0)
	for _, m := range metrics {
		tree.add(m)
	}
	return tree.metricGroups()
}

func Generate(name, datasource, selector string, filerVariables []string, metrics map[string]prometheus.Metric) (*dashboard.Builder, error) {
	metricGroups := groupMetrics(metrics)
	dashboard, err := newDashboard(name, datasource, selector, filerVariables, metricGroups)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}
