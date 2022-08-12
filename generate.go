package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fusakla/autograf/packages/generator"
	"github.com/fusakla/autograf/packages/grafana"
	"github.com/fusakla/autograf/packages/prometheus"
)

func openInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	return err
}

func (r *Command) Run(ctx *Context) error {
	var metrics map[string]prometheus.Metric
	if r.MetricsFile != "" {
		data, err := os.ReadFile(r.MetricsFile)
		if err != nil {
			return err
		}
		metrics, err = prometheus.ParseMetricsText(data, r.OpenMetricsFormat)
		if err != nil {
			return err
		}
	} else if r.PrometheusURL != "" {
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
	dashboard, err := generator.Generate(r.GrafanaDashboardName, r.GrafanaDataSource, r.Selector, r.GrafanaVariables, metrics)
	if err != nil {
		return err
	}
	if r.GrafanaURL != "" {
		if r.grafanaToken == "" {
			return fmt.Errorf("you have to specify the GRAFANA_TOKEN variable")
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		dashboardUrl, err := grafana.UpsertDashboard(ctx, r.GrafanaURL, r.grafanaToken, r.GrafanaFolder, dashboard)
		if err != nil {
			return err
		}
		fmt.Println("Dashboard successfully generated, see " + dashboardUrl)
		if r.OpenBrowser {
			if err := openInBrowser(dashboardUrl); err != nil  {
				return err
			}
		}
	} else {
		jsonData, err := dashboard.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	}
	return nil
}
