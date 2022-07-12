package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

var (
	Version = "development"
)

type Context struct {
	logger logrus.FieldLogger
}

type VersionCmd struct{}

func (r *VersionCmd) Run(_ *Context) error {
	fmt.Println(Version)
	return nil
}

var help = `
Autograf generates Grafana dashboard form Prometheus metrics either read form a /metrics endpoint or queried form live Prometheus instance.
The dashboard JSON is by default printed to stdout. But can also upload the dashboard directly to your Grafana instance.
You can configure most of the flags using config file. See the docs.

Example from /metrics:
  curl http://foo.bar/metrics | autograf --metrics-file - 

Example from Prometheus query:
  GRAFANA_TOKEN=xxx autograf --prometheus-url http://prometheus.foo --selector {app='foo'} --grafana-url http://grafana.bar 

`

type Command struct {
	Debug bool `help:"Enable debug logging"`

	MetricsFile       string `help:"File containing the metrics exposed by app (will read stdin if se to - )"`
	OpenMetricsFormat bool   `help:"Metrics data are in the application/openmetrics-text format."`

	PrometheusURL    string   `help:"URL of Prometheus instance to fetch the metrics from."`
	Selector         string   `help:"Selector to filter metrics from the Prometheus instance."`
	GrafanaVariables []string `help:"Labels used as a variables for filtering in dashboard"`

	GrafanaURL           string `help:"URL of Grafana to upload the dashboard to, if not specified, dashboard JSON is printed to stdout"`
	GrafanaFolder        string `help:"Name of target Grafana folder"`
	GrafanaDashboardName string `help:"Name of the Grafana dashboard"`
	GrafanaDataSource    string `help:"Name of the Grafana datasource to use"`

	OpenBrowser bool `help:"Open the Grafana dashboard automatically in browser"`

	grafanaToken string `kong:"-"`
}

func (c *Command) updateFromConfig(conf config) {
	if c.PrometheusURL == "" {
		c.PrometheusURL = conf.PrometheusURL
	}
	if c.GrafanaURL == "" {
		c.GrafanaURL = conf.GrafanaURL
	}
	if c.GrafanaFolder == "" {
		c.GrafanaFolder = conf.GrafanaFolder
	}
	if c.GrafanaDashboardName == "" {
		c.GrafanaDashboardName = conf.GrafanaDashboardName
		if c.GrafanaDashboardName == "" {
			c.GrafanaDashboardName = "Autograf dashboard"
		}
	}
	if c.GrafanaDataSource == "" {
		c.GrafanaDataSource = conf.GrafanaDataSource
	}
	if c.grafanaToken == "" {
		c.grafanaToken = conf.GrafanaToken
	}
	if !c.OpenBrowser && conf.OpenBrowser {
		c.OpenBrowser = true
	}
	if len(c.GrafanaVariables) == 0 {
		c.GrafanaVariables = conf.GrafanaVariables
		if len(c.GrafanaVariables) == 0 {
			c.GrafanaVariables = []string{"job", "instance"}
		}
	}
}

var CLI Command

func main() {
	ctx := kong.Parse(&CLI, kong.Description(help))
	rootLogger := logrus.New()
	rootLogger.SetOutput(os.Stderr)
	rootLogger.SetLevel(logrus.WarnLevel)
	if CLI.Debug {
		rootLogger.SetLevel(logrus.DebugLevel)
	}
	CLI.grafanaToken = os.Getenv("GRAFANA_TOKEN")
	CLI.updateFromConfig(loadConfig(rootLogger))

	if CLI.PrometheusURL == "" && CLI.MetricsFile == "" {
		rootLogger.Error("Error, at leas one of the --prometheus-url or --metrics-file flags have to be set")
		os.Exit(1)
	}

	err := ctx.Run(&Context{
		logger: rootLogger.WithField("command", ctx.Command()),
	})
	ctx.FatalIfErrorf(err)
}
