package generator

import (
	"fmt"
	"strings"

	"github.com/fusakla/autograf/packages/model"
	"github.com/fusakla/sdk"
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
	metric *model.Metric
	level  int
}

func (t *metricsTree) metrics() []*model.Metric {
	metrics := []*model.Metric{}
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

func (t *metricsTree) metricGroups() map[string][]*model.Metric {
	others := []*model.Metric{}
	groups := map[string][]*model.Metric{}
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

func (t *metricsTree) add(metric *model.Metric) {
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

func groupMetrics(metrics map[string]*model.Metric) map[string][]*model.Metric {
	tree := newTree("", 0)
	metricGroups := map[string][]*model.Metric{}
	for _, m := range metrics {
		if m.Config.Row != "" {
			r, ok := metricGroups[m.Config.Row]
			if !ok {
				metricGroups[m.Config.Row] = []*model.Metric{m}
			} else {
				metricGroups[m.Config.Row] = append(r, m)
			}
			continue
		}
		tree.add(m)
	}
	for k, v := range tree.metricGroups() {
		metricGroups[k] = v
	}
	return metricGroups
}

func Generate(name, datasource, selector string, filerVariables []string, metrics map[string]*model.Metric) (*sdk.Board, error) {
	metricGroups := groupMetrics(metrics)	
	dashboard, err := newDashboard(name, datasource, selector, filerVariables, metricGroups)
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}
