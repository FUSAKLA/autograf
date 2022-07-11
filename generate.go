package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/fusakla/autograf/packages/generator"
	"github.com/fusakla/autograf/packages/grafana"
	"github.com/fusakla/autograf/packages/prometheus"
)

type GenerateCmd struct {
	MetricsFile       *os.File `group:"input" xor:"fromm-metrics" required:"" help:"File containing the metrics exposed by app (will read stdin if se to -)"`
	OpenMetricsFormat bool     `help:"Metrics data are in the application/openmetrics-text format."`

	PrometheusURL *url.URL `group:"input" xor:"fromm-metrics" required:"" help:"URL of Prometheus instance to fetch the metrics from."`
	Selector      string   `default:"{}" help:"Selector to filter metrics from the Prometheus instance."`
	AggregateBy   []string `help:"Additional labels to aggregate the queries by."`

	DashboardName string `default:"Autograf dashboard" help:"Name of the dashboard"`
	DataSource    string `help:"Name of the data source to use"`

	GrafanaURL    *url.URL `help:"URL of Grafana to upload the dashboard to, if not specified, dashboard JSON is printed to stdout"`
	GrafanaFolder string   `help:"Name of target Grafana folder"`
}

func (r *GenerateCmd) Run(ctx *Context) error {
	var metrics map[string]prometheus.Metric
	if r.MetricsFile != nil {
		data, err := io.ReadAll(r.MetricsFile)
		if err != nil {
			return err
		}
		metrics, err = prometheus.ParseMetricsText(data, r.OpenMetricsFormat)
		if err != nil {
			return err
		}
	} else if r.PrometheusURL != nil {
		client, err := prometheus.NewClient(ctx.logger, r.PrometheusURL, nil)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		metrics, err = client.MetricsForSelector(ctx, r.Selector)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("At least one of inputs metrics file or Prometheus URL is required")
	}
	dashboard, err := generator.Generate(r.DashboardName, r.DataSource, r.Selector, metrics)
	if err != nil {
		return err
	}
	if r.GrafanaURL != nil {
		gToken, ok := os.LookupEnv("GRAFANA_TOKEN")
		if !ok {
			return fmt.Errorf("you have to specify the GRAFANA_TOKEN variable")
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		dashboardUrl, err := grafana.UpsertDashboard(ctx, r.GrafanaURL, gToken, r.GrafanaFolder, dashboard)
		if err != nil {
			return err
		}
		fmt.Println("Dashboard successfully generated, see " + dashboardUrl)
	} else {
		jsonData, err := dashboard.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	}
	return nil
}
