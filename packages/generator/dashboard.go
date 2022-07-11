package generator

import (
	"fmt"
	"sort"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/variable/datasource"
	"github.com/fusakla/autograf/packages/prometheus"
)

func newRow(dataSource string, selector string, name string, metrics []*prometheus.Metric) dashboard.Option {
	panels := []row.Option{row.Collapse()}
	for _, m := range metrics {
		panels = append(panels, newPanel(dataSource, selector, *m))
	}
	return dashboard.Row(
		name,
		panels...,
	)
}

func newDashboard(name, datasourceName, selector string, metricGroups map[string][]*prometheus.Metric) (dashboard.Builder, error) {
	opts := []dashboard.Option{
		dashboard.AutoRefresh("1m"),
		dashboard.SharedCrossHair(),
		dashboard.Tags([]string{"autograf", "generated"}),
		dashboard.Timezone(dashboard.Browser),
		dashboard.Time("now-1h", "now"),
		dashboard.VariableAsDatasource("datasource", datasource.Type("prometheus"), datasource.Regex(fmt.Sprintf("/%s/", datasourceName))),
	}
	rowNames := []string{}
	for k := range metricGroups {
		rowNames = append(rowNames, k)
	}
	sort.Strings(rowNames)
	for _, r := range rowNames {
		opts = append(opts, newRow("${datasource}", selector, r, metricGroups[r]))
	}
	return dashboard.New(
		name,
		opts...,
	)
}
