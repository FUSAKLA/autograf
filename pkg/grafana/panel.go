package grafana

import (
	"fmt"
	"regexp"

	"github.com/fusakla/autograf/pkg/generator"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/heatmap"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/table"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

const (
	timeSeriesFormat      = "time_series"
	panelHeightCoeficient = 1
	rateIntervalVariable  = "$__rate_interval"
)

func panelNameFromQuery(query string) string {
	return regexp.MustCompile(`\{[^{}]*\}`).ReplaceAllString(query, "")
}

func addLimitTarget(panel *timeseries.PanelBuilder, lType generator.LimitType, metric string, selector string) {
	panel.WithTarget(prometheus.NewDataqueryBuilder().
		RefId(metric).
		Expr(generator.ThresholdQuery(metric, selector, lType)).
		Range().
		LegendFormat(fmt.Sprintf("%s limit", lType)).
		Format(prometheus.PromQueryFormatTimeSeries),
	).OverrideByName(metric, []dashboard.DynamicConfigValue{
		{Id: "custom.fillOpacity", Value: 0},
		{Id: "color", Value: map[string]string{"mode": "fixed", "fixedColor": "red"}},
		{Id: "custom.lineWidth", Value: 3},
		{Id: "custom.lineStyle", Value: map[string]string{"fill": "dash"}},
	})
}

func newTimeSeriesPanel(dataSource dashboard.DataSourceRef, selector string, metric generator.Metric) cog.Builder[dashboard.Panel] {
	query := metric.PromQlQuery(selector, rateIntervalVariable)
	panel := timeseries.NewPanelBuilder().
		Title(panelNameFromQuery(query)).
		Description(metric.Help).
		Datasource(dataSource).
		Legend(common.NewVizLegendOptionsBuilder().
			DisplayMode(common.LegendDisplayModeTable).
			Calcs(metric.Config.LegendCalcs).
			ShowLegend(false)).
		Unit(string(metric.Unit)).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeSingle).
			Sort(common.SortOrderDescending),
		).LineWidth(float64(metric.Config.LineWidth)).
		LineStyle(common.NewLineStyleBuilder().
			Fill("solid"),
		).DrawStyle("line").
		ShowPoints("auto").
		PointSize(1).
		Span(uint32(metric.Config.Width)).
		Height(uint32(metric.Config.Height * panelHeightCoeficient))

	if metric.Config.Stack {
		panel.Stacking(common.NewStackingConfigBuilder().Mode("normal"))
	}

	switch metric.Config.Scale {
	case "linear":
		panel.ScaleDistribution(common.NewScaleDistributionConfigBuilder().Type(common.ScaleDistributionLinear))
	case "log2":
		panel.ScaleDistribution(common.NewScaleDistributionConfigBuilder().Type(common.ScaleDistributionLog).Log(2))
	case "log10":
		panel.ScaleDistribution(common.NewScaleDistributionConfigBuilder().Type(common.ScaleDistributionLog).Log(10))
	}

	panel.WithTarget(prometheus.NewDataqueryBuilder().Datasource(dataSource).
		RefId(metric.Name).
		Expr(query).
		Range().
		Format(prometheus.PromQueryFormatTimeSeries),
	)

	if metric.Config.MaxFromMetric != "" {
		addLimitTarget(panel, generator.LimitMax, metric.Config.MaxFromMetric, selector)
	}
	if metric.Config.MinFromMetric != "" {
		addLimitTarget(panel, generator.LimitMin, metric.Config.MinFromMetric, selector)
	}

	return panel
}

func newHeatmapPanel(dataSource dashboard.DataSourceRef, selector string, metric generator.Metric) cog.Builder[dashboard.Panel] {
	query := metric.PromQlQuery(selector, rateIntervalVariable)

	panel := heatmap.NewPanelBuilder().
		Title(panelNameFromQuery(query)).
		Description(metric.Help).
		Unit(string(metric.Unit)).
		Datasource(dataSource).
		ShowLegend().
		Calculate(false).
		Color(heatmap.NewHeatmapColorOptionsBuilder().
			Mode(heatmap.HeatmapColorModeOpacity).
			Exponent(0.3).
			Steps(20).
			Fill("super-light-blue"),
		).
		ShowColorScale(true).
		ShowYHistogram().
		YAxis(heatmap.NewYAxisConfigBuilder().
			AxisPlacement("left").
			Unit(string(metric.Unit)),
		).
		CellGap(1).
		CellValues(heatmap.NewCellValuesBuilder().
			Unit(string(metric.Unit)),
		).
		Span(uint32(metric.Config.Width)).
		Height(uint32(metric.Config.Height * panelHeightCoeficient))

	panel.WithTarget(prometheus.NewDataqueryBuilder().
		Datasource(dataSource).
		RefId(metric.Name).
		Expr(query).
		Range().
		Format(prometheus.PromQueryFormatHeatmap).
		LegendFormat("{{le}}"),
	)

	return panel
}

func newInfoPanel(dataSource dashboard.DataSourceRef, selector string, metric generator.Metric) cog.Builder[dashboard.Panel] {
	panel := table.NewPanelBuilder().
		DisplayName(panelNameFromQuery(metric.Name)).
		Description(metric.Help).
		OverrideByRegexp(
			"(__name__|Time|Value)",
			[]dashboard.DynamicConfigValue{{Id: "custom.hidden", Value: true}},
		).
		Span(24).
		Height(uint32(metric.Config.Height * panelHeightCoeficient)).
		WithTarget(prometheus.NewDataqueryBuilder().
			Datasource(dataSource).
			RefId(metric.Name).
			Expr(metric.PromQlQuery(selector, rateIntervalVariable)).
			Instant().
			Format(prometheus.PromQueryFormatTable),
		)
	return panel
}

func newPanel(dataSource dashboard.DataSourceRef, selector string, metric generator.Metric) (cog.Builder[dashboard.Panel], bool) {
	switch metric.MetricType {
	case "gauge":
		return newTimeSeriesPanel(dataSource, selector, metric), false
	case "counter":
		return newTimeSeriesPanel(dataSource, selector, metric), false
	case "summary":
		return newTimeSeriesPanel(dataSource, selector, metric), false
	case "histogram":
		return newHeatmapPanel(dataSource, selector, metric), false
	case "info":
		return newInfoPanel(dataSource, selector, metric), true
	}
	metric.Help = fmt.Sprintf("WARNING: Unknown metric type %s!\n\n%s", metric.MetricType, metric.Help)
	return newTimeSeriesPanel(dataSource, selector, metric), false
}
