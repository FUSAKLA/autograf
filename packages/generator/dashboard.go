package generator

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/variable/datasource"
	"github.com/K-Phoen/grabana/variable/query"
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

func selectorWithVariablesFilter(selector string, filerVariables []string) string {
	new := strings.TrimSuffix(selector, "}")
	for _, v := range filerVariables {
		new += fmt.Sprintf(",%s=~'$%s'", v, v)
	}
	return new + "}"
}

func labelVariable(datasourceName, selector, name string) dashboard.Option {
	return dashboard.VariableAsQuery(
		name,
		query.DataSource(datasourceName),
		query.Request(fmt.Sprintf("query_result(%s)", selector)),
		query.Regex(fmt.Sprintf(`/%s="([^"]+)"/`, name)),
		query.AllValue(".*"),
		query.DefaultAll(),
		query.IncludeAll(),
		query.Multi(),
		query.Sort(query.AlphabeticalAsc), query.DefaultAll(),
	)
}

func newDashboard(name, datasourceName, selector string, filerVariables []string, metricGroups map[string][]*prometheus.Metric) (dashboard.Builder, error) {
	opts := []dashboard.Option{
		dashboard.AutoRefresh("1m"),
		dashboard.SharedCrossHair(),
		dashboard.Tags([]string{"autograf", "generated"}),
		dashboard.Timezone(dashboard.Browser),
		dashboard.Time("now-1h", "now"),
		dashboard.VariableAsDatasource(
			"datasource",
			datasource.Type("prometheus"),
			datasource.Regex(fmt.Sprintf("/%s/", regexp.QuoteMeta(datasourceName))),
		),
	}
	for _, v := range filerVariables {
		opts = append(opts, labelVariable(datasourceName, selector, v))
	}
	rowNames := []string{}
	for k := range metricGroups {
		rowNames = append(rowNames, k)
	}
	sort.Strings(rowNames)
	for _, r := range rowNames {
		opts = append(opts, newRow("${datasource}", selectorWithVariablesFilter(selector, filerVariables), r, metricGroups[r]))
	}
	return dashboard.New(
		name,
		opts...,
	)
}
