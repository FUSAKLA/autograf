package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type config struct {
	PrometheusURL string `json:"prometheus_url,omitempty"`

	GrafanaURL           string   `json:"grafana_url,omitempty"`
	GrafanaToken         string   `json:"grafana_token,omitempty"`
	GrafanaFolder        string   `json:"grafana_folder,omitempty"`
	GrafanaDashboardName string   `json:"grafana_dashboard_name,omitempty"`
	GrafanaDataSource    string   `json:"grafana_datasource,omitempty"`
	GrafanaVariables     []string `json:"grafana_variables,omitempty"`

	OpenBrowser bool `json:"open_browser,omitempty"`
}

func loadConfig(logger logrus.FieldLogger) config {
	var c config
	configFilePath, ok := os.LookupEnv("AUTOGRAF_CONFIG")
	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.WithField("error", err).Warn("failed to load autograf config from home dir")
			return c
		}
		configFilePath = path.Join(home, ".autograf.json")
	}
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return c
	}
	if err := json.Unmarshal(data, &c); err != nil {
		logger.WithField("error", err).Warn("invalid config file format")
	}
	return c
}
