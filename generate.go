package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fusakla/autograf/packages/generator"
	"github.com/fusakla/autograf/packages/grafana"
	"github.com/fusakla/autograf/packages/model"
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
	var metrics map[string]*model.Metric
	if r.MetricsFile != "" {
		var data []byte
		var err error
		if r.MetricsFile == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(strings.TrimSpace(r.MetricsFile))
		}
		if err != nil {
			return err
		}
		metrics, err = prometheus.ParseMetricsText(data, r.OpenMetricsFormat)
		if err != nil {
			return err
		}
	} else if r.PrometheusURL != "" {
		client, err := prometheus.NewClient(ctx.logger, strings.TrimSpace(r.PrometheusURL), nil)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		metrics, err = client.MetricsForSelector(ctx, strings.TrimSpace(r.Selector))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("At least one of inputs metrics file or Prometheus URL is required")
	}
	if err := prometheus.ProcessMetrics(metrics); err != nil {
		return err
	}
	dashboard, err := generator.Generate(strings.TrimSpace(r.GrafanaDashboardName), strings.TrimSpace(r.GrafanaDataSource), strings.TrimSpace(r.Selector), r.GrafanaVariables, metrics)
	if err != nil {
		return err
	}
	if r.GrafanaURL != "" {
		if r.grafanaToken == "" {
			return fmt.Errorf("you have to specify the GRAFANA_TOKEN variable")
		}
		cli := grafana.NewClient(r.GrafanaURL, r.grafanaToken)
		folderUid, err := cli.EnsureFolder(strings.TrimSpace(r.GrafanaFolder))
		if err != nil {
			return err
		}
		dashboardUrl, err := cli.UpsertDashboard(folderUid, dashboard)
		if err != nil {
			return err
		}
		fmt.Println("Dashboard successfully generated, see " + dashboardUrl)
		if r.OpenBrowser {
			if err := openInBrowser(dashboardUrl); err != nil {
				return err
			}
		}
	} else {
		jsonData, err := json.Marshal(dashboard)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	}
	return nil
}
