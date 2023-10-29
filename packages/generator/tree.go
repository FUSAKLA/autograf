package generator

import (
	"fmt"
	"strings"
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
	metric *Metric
	level  int
}

func (t *metricsTree) metrics() []*Metric {
	metrics := []*Metric{}
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

func (t *metricsTree) metricGroups() map[string][]*Metric {
	others := []*Metric{}
	groups := map[string][]*Metric{}
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
	if len(others) == 1 {
		groups[others[0].Name] = others
	} else if len(others) > 0 {
		groups[t.prefix] = others
	}
	return groups
}

func (t metricsTree) String() string {
	indent := strings.Repeat(" ", t.level)
	out := ""
	out += fmt.Sprintf(indent + t.prefix + "\n")
	if t.metric != nil {
		out += fmt.Sprintf("%s- %v", indent, t.metric)
	}
	for _, l := range t.leafs {
		out += "\n" + l.String()
	}
	return out
}

func (t *metricsTree) add(metric *Metric) {
	strippedMetricName := strings.TrimPrefix(strings.TrimPrefix(metric.Name, t.prefix), "_")
	if strippedMetricName == "" {
		t.metric = metric
		return
	}
	parts := strings.SplitN(strippedMetricName, "_", 2)
	prefix := parts[0]
	if _, ok := t.leafs[prefix]; !ok {
		t.leafs[prefix] = newTree(strings.TrimPrefix(t.prefix+"_"+prefix, "_"), t.level+1)
	}
	t.leafs[prefix].add(metric)
}

func groupIntoPseudoDashboard(metrics map[string]*Metric) PseudoDashboard {
	tree := newTree("", 0)
	dashboard := PseudoDashboard{}
	for _, m := range metrics {
		if m.Config.Row != "" {
			dashboard.AddRowPanels(m.Config.Row, []*Metric{m})
			continue
		}
		tree.add(m)
	}
	for k, v := range tree.metricGroups() {
		dashboard.AddRowPanels(k, v)
	}
	return dashboard
}
