package grafana

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fusakla/autograf/pkg/generator"
	"github.com/fusakla/sdk"
	"golang.org/x/exp/maps"
)

func newRow(dataSource *sdk.DatasourceRef, selector string, name string, metrics []*generator.Metric) *sdk.Row {
	row := sdk.Row{
		Title:    name,
		Collapse: true,
	}
	metricNames := []string{}
	metricNamesChars := 0
	trimMetricNames := false
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
		if m.MetricType == generator.MetricTypeInfo {
			row.Panels = append([]sdk.Panel{*newPanel(dataSource, selector, *m)}, row.Panels...)
		} else {
			row.Panels = append(row.Panels, *newPanel(dataSource, selector, *m))
		}

	}
	if len(metricNames) > 1 {
		row.Title += " ❯ " + strings.Join(metricNames, " ❙ ")
		if trimMetricNames {
			row.Title += " | ..."
		}
	}
	return &row
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

func labelVariable(datasourceName *sdk.DatasourceRef, selector, name string) sdk.TemplateVar {
	if selector == "" {
		selector = "up"
	}
	return sdk.TemplateVar{
		Name: name,
		Type: "query",
		Datasource: &sdk.DatasourceRef{
			Type: "prometheus",
			UID:  "${datasource}",
		},
		Query: struct {
			Query string `json:"query"`
			RefId string `json:"refId"`
		}{
			Query: fmt.Sprintf("query_result(%s)", selector),
			RefId: "StandardVariableQuery",
		},
		Regex:      fmt.Sprintf(`/%s="([^"]+)"/`, name),
		AllValue:   ".*",
		Options:    []sdk.Option{},
		IncludeAll: true,
		Multi:      true,
		Sort:       1, // Alphabetical ASC
	}
}

func NewDashboard(name, datasourceName, selector string, filterVariables []string, pseudoDashboard generator.PseudoDashboard) (*sdk.Board, error) {
	board := sdk.NewBoard(name)
	board.Refresh = &sdk.BoolString{Flag: true, Value: "1m"}
	board.GraphTooltip = 1 // 0 for no shared crosshair or tooltip (default), 1 for shared crosshair, 2 for shared crosshair AND shared tooltip
	board.Tags = []string{"autograf", "generated"}
	board.Time = sdk.Time{From: "now-1h", To: "now"}
	board.Timezone = "browser"
	var refresh int64 = 1
	board.Templating.List = []sdk.TemplateVar{
		{
			Type:    "datasource",
			Name:    "datasource",
			Label:   "Datasource",
			Query:   "prometheus",
			Refresh: sdk.BoolInt{Flag: true, Value: &refresh},
			Options: []sdk.Option{},
			Current: sdk.Current{
				Text:     &sdk.StringSliceString{Valid: true, Value: []string{datasourceName}},
				Selected: true,
				Value:    datasourceName,
			},
		},
	}
	for _, v := range filterVariables {
		board.Templating.List = append(board.Templating.List, labelVariable(&sdk.DatasourceRef{UID: "${datasource}"}, selector, v))
	}
	rowNames := maps.Keys(pseudoDashboard.Rows)
	slices.Sort(rowNames)
	for _, r := range rowNames {
		board.Rows = append(board.Rows, newRow(&sdk.DatasourceRef{Type: "prometheus", UID: "${datasource}"}, selectorWithVariablesFilter(selector, filterVariables), r, pseudoDashboard.Rows[r].Metrics))
	}
	return board, nil
}
