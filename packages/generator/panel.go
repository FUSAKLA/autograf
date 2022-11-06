package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fusakla/autograf/packages/model"
	"github.com/fusakla/sdk"
)

var timeMetricsRegexp = regexp.MustCompile(`.+(_time|_time_seconds|_timestamp|_timestamp_seconds|_update)$`)

type limitType string

const (
	maxLimit limitType = "max"
	minLimit limitType = "min"

	timeSeriesFormat = "time_series"

	panelHeightCoeficient = 40

	rateIntervalVariable = "$__rate_interval"
)

func panelNameFromQuery(query string) string {
	return regexp.MustCompile(`\{[^{}]*\}`).ReplaceAllString(query, "")
}

func addLimitTarget(panel *sdk.Panel, lType limitType, metric string, selector string) {
	panel.TimeseriesPanel.Targets = append(panel.TimeseriesPanel.Targets, sdk.Target{
		RefID:        metric,
		Expr:         fmt.Sprintf("%s(%s%s)", lType, metric, selector),
		Instant:      false,
		LegendFormat: fmt.Sprintf("%s limit", lType),
		Format:       timeSeriesFormat,
	})
	panel.TimeseriesPanel.FieldConfig.Overrides = append(panel.TimeseriesPanel.FieldConfig.Overrides, sdk.FieldConfigOverride{
		Properties: []sdk.FieldConfigOverrideProperty{
			{ID: "custom.fillOpacity", Value: 0},
			{ID: "color", Value: map[string]string{"mode": "fixed", "fixedColor": "red"}},
			{ID: "custom.lineWidth", Value: 3},
			{ID: "custom.lineStyle", Value: map[string]string{"fill": "dash"}},
		},
	})
}

func newTimeSeriesPanel(dataSource *sdk.DatasourceRef, selector string, metric model.Metric) *sdk.Panel {
	query := fmt.Sprintf("%s%s", metric.Name, selector)
	if metric.MetricType == model.MetricTypeCounter {
		query = fmt.Sprintf("rate(%s%s[%s])", metric.Name, selector, rateIntervalVariable)
	}
	if timeMetricsRegexp.MatchString(metric.Name) {
		query = "time() - " + query
		metric.Unit = model.MetricUnitSeconds
	}
	if metric.Config.Aggregation != "" {
		query = fmt.Sprintf("%s by (%s) (%s)", metric.Config.Aggregation, strings.Join(metric.Config.AggregateBy, ","), query)
	}

	panel := sdk.NewTimeseries(panelNameFromQuery(query))
	panel.Description = &metric.Help
	panel.Datasource = &sdk.DatasourceRef{
		LegacyName: "$datasource",
	}
	panel.TimeseriesPanel.Options.Legend.ShowLegend = false
	panel.TimeseriesPanel.Options.Legend.Calcs = metric.Config.LegendCalcs
	panel.TimeseriesPanel.Options.Legend.DisplayMode = "table"
	panel.TimeseriesPanel.FieldConfig.Defaults.Unit = string(metric.Unit)
	panel.TimeseriesPanel.Options.Tooltip.Mode = "single"
	panel.TimeseriesPanel.Options.Tooltip.Sort = "desc"
	panel.TimeseriesPanel.FieldConfig.Defaults.Custom.LineWidth = metric.Config.LineWidth
	panel.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle = "line"
	panel.TimeseriesPanel.FieldConfig.Defaults.Custom.LineStyle.Fill = "solid"
	panel.TimeseriesPanel.FieldConfig.Defaults.Custom.ShowPoints = "auto"
	panel.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize = 1
	if metric.Config.Stack {
		panel.TimeseriesPanel.FieldConfig.Defaults.Custom.Stacking.Mode = "normal"
	}
	panel.Span = metric.Config.Width
	panel.Height = metric.Config.Height * panelHeightCoeficient

	switch metric.Config.Scale {
	case "linear":
		panel.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Type = "linear"
	case "log2":
		panel.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Type = "log"
		panel.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Log = 2
	case "log10":
		panel.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Type = "log"
		panel.TimeseriesPanel.FieldConfig.Defaults.Custom.ScaleDistribution.Log = 10
	}

	panel.TimeseriesPanel.Targets = append(panel.TimeseriesPanel.Targets, sdk.Target{
		Datasource: dataSource,
		RefID:      metric.Name,
		Expr:       query,
		Instant:    false,
		Format:     "time_series",
	})

	if metric.Config.MaxFromMetric != "" {
		addLimitTarget(panel, maxLimit, metric.Config.MaxFromMetric, selector)
	}
	if metric.Config.MinFromMetric != "" {
		addLimitTarget(panel, minLimit, metric.Config.MinFromMetric, selector)
	}

	return panel
}

func newHeatmapPanel(dataSource *sdk.DatasourceRef, selector string, metric model.Metric) *sdk.Panel {
	query := fmt.Sprintf("sum(rate(%s%s[%s])) by (%s)", metric.Name, selector, strings.Join(append(metric.Config.AggregateBy, "le"), ","), rateIntervalVariable)
	panel := sdk.NewHeatmap(panelNameFromQuery(query))
	panel.Description = &metric.Help
	panel.HeatmapPanel.HideZeroBuckets = true
	panel.HeatmapPanel.DataFormat = "tsbuckets"
	panel.HeatmapPanel.FieldConfig.Defaults.Unit = string(metric.Unit)
	panel.HeatmapPanel.Options.Tooltip.Show = true
	panel.HeatmapPanel.Options.Tooltip.ShowHistogram = true
	panel.HeatmapPanel.Options.Calculate = true
	panel.HeatmapPanel.Options.YAxis.AxisPlacement = "left"
	panel.HeatmapPanel.Options.YAxis.Unit = string(metric.Unit)
	panel.HeatmapPanel.Options.Color.Mode = "opacity"
	panel.HeatmapPanel.Options.Color.Exponent = 0.3
	panel.HeatmapPanel.Options.Color.Fill = "super-light-blue"
	panel.HeatmapPanel.Options.CellGap = 1
	panel.HeatmapPanel.Options.Legend.Show = true
	panel.HeatmapPanel.CellGap = 1
	panel.HeatmapPanel.CellValues.Unit = "number"
	panel.Span = metric.Config.Width
	panel.Height = metric.Config.Height * panelHeightCoeficient
	panel.HeatmapPanel.Targets = append(panel.HeatmapPanel.Targets, sdk.Target{
		Datasource:   dataSource,
		RefID:        metric.Name,
		Expr:         query,
		Instant:      false,
		Format:       "heatmap",
		LegendFormat: "{{le}}",
	})
	return panel
}

func newInfoPanel(dataSource *sdk.DatasourceRef, selector string, metric model.Metric) *sdk.Panel {
	panel := sdk.NewTable(metric.Name)
	panel.Description = &metric.Help
	panel.TablePanel.FieldConfig.Overrides = []sdk.FieldConfigOverride{
		{
			Matcher: struct {
				ID      string `json:"id"`
				Options string `json:"options"`
			}{ID: "byRegexp", Options: "(__name__|Time|Value)"},
			Properties: []sdk.FieldConfigOverrideProperty{
				{ID: "custom.hidden", Value: "true"},
			},
		}}
	panel.Span = 12
	panel.Height = metric.Config.Height * panelHeightCoeficient
	panel.TablePanel.Targets = append(panel.TablePanel.Targets, sdk.Target{
		Datasource: dataSource,
		RefID:      metric.Name,
		Expr:       fmt.Sprintf("%s%s", metric.Name, selector),
		Instant:    true,
		Format:     "table",
	})
	return panel
}

func newPanel(dataSource *sdk.DatasourceRef, selector string, metric model.Metric) *sdk.Panel {
	switch metric.MetricType {
	case "gauge":
		return newTimeSeriesPanel(dataSource, selector, metric)
	case "counter":
		return newTimeSeriesPanel(dataSource, selector, metric)
	case "summary":
		return newTimeSeriesPanel(dataSource, selector, metric)
	case "histogram":
		return newHeatmapPanel(dataSource, selector, metric)
	case "info":
		return newInfoPanel(dataSource, selector, metric)
	}
	metric.Help = fmt.Sprintf("WARNING: Unknown metric type %s!\n\n%s", metric.MetricType, metric.Help)
	return newTimeSeriesPanel(dataSource, selector, metric)
}
