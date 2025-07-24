package grafana

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fusakla/autograf/pkg/generator"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func newRow(dataSource dashboard.DataSourceRef, selector string, name string, metrics []*generator.Metric) cog.Builder[dashboard.RowPanel] {
	row := dashboard.NewRowBuilder(name).Collapsed(true)
	metricNames := []string{}
	metricNamesChars := 0
	trimMetricNames := false
	panels := []cog.Builder[dashboard.Panel]{}
	for _, m := range metrics {
		mName := strings.TrimPrefix(m.Name, name)
		if metricNamesChars+len(mName) < 150 {
			metricNames = append(metricNames, mName)
			metricNamesChars += len(mName)
		} else {
			trimMetricNames = true
		}
		if len(metrics) == 1 {
			m.Config.Width = 12
		}
		panel, isInfo := newPanel(dataSource, selector, *m)
		if isInfo {
			panels = append([]cog.Builder[dashboard.Panel]{panel}, panels...)
		} else {
			panels = append(panels, panel)
		}
	}
	for _, p := range panels {
		row = row.WithPanel(p)
	}
	if len(metricNames) > 1 {
		newTitle := name + " ❯ " + strings.Join(metricNames, " ❙ ")
		if trimMetricNames {
			newTitle += " | ..."
		}
		row = row.Title(newTitle)
	}
	return row
}

func selectorWithVariablesFilter(selector string, filerVariables []string) string {
	new := strings.TrimSuffix(selector, "}") + ","
	if selector == "" {
		new = "{"
	}
	filters := make([]string, len(filerVariables))
	for i, v := range filerVariables {
		filters[i] = fmt.Sprintf("%s=~'$%s'", v, v)
	}
	return new + strings.Join(filters, ",") + "}"
}

func labelVariable(datasource dashboard.DataSourceRef, selector, name string) *dashboard.QueryVariableBuilder {
	if selector == "" {
		selector = "up"
	}
	return dashboard.NewQueryVariableBuilder(name).
		Datasource(datasource).
		Query(dashboard.StringOrMap{String: cog.ToPtr(fmt.Sprintf("label_values(%s, %s)", selector, name))}).
		// Regex(fmt.Sprintf(`/%s="([^"]+)"/`, name)).
		AllValue(".*").
		IncludeAll(true).
		Multi(true).
		Sort(dashboard.VariableSortAlphabeticalAsc).
		Refresh(dashboard.VariableRefreshOnTimeRangeChanged)
}

func NewDashboard(name, datasourceID, selector string, filterVariables []string, pseudoDashboard generator.PseudoDashboard) *dashboard.DashboardBuilder {
	board := dashboard.NewDashboardBuilder(name).
		Refresh("1m").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Tags([]string{"autograf", "generated"}).
		Time("now-1h", "now").
		Timezone("browser").
		WithVariable(
			dashboard.NewDatasourceVariableBuilder("datasource").
				Label("Datasource").
				Type("prometheus").
				Current(dashboard.VariableOption{
					Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr(datasourceID)},
					Selected: cog.ToPtr(true),
				}),
		)
	datasource := dashboard.DataSourceRef{Type: cog.ToPtr("prometheus"), Uid: cog.ToPtr("${datasource}")}

	for _, v := range filterVariables {
		board = board.WithVariable(labelVariable(datasource, selector, v))
	}
	rowNames := make([]string, 0, len(pseudoDashboard.Rows))
	for k := range pseudoDashboard.Rows {
		rowNames = append(rowNames, k)
	}
	slices.Sort(rowNames)
	for _, r := range rowNames {
		board = board.WithRow(newRow(datasource, selectorWithVariablesFilter(selector, filterVariables), r, pseudoDashboard.Rows[r].Metrics))
	}
	return board
}
