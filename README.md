# Autograf
[![Go Report
Card](https://goreportcard.com/badge/github.com/fusakla/autograf)](https://goreportcard.com/report/github.com/fusakla/autograf)
[![GitHub actions
CI](https://img.shields.io/github/actions/workflow/status/fusakla/autograf/go.yaml)](https://github.com/fusakla/autograf/actions?query=branch%3Amain)
[![Docker Pulls](https://img.shields.io/docker/pulls/fusakla/autograf)](https://hub.docker.com/r/fusakla/autograf)
[![GitHub binaries
download](https://img.shields.io/github/downloads/fusakla/autograf/total?label=Prebuilt%20binaries%20downloads)](https://github.com/fusakla/autograf/releases/latest)

**Dynamically generate Grafana dashboard based on Prometheus metrics**

<p align="center"><img src="./autograf.excalidraw.png"></p>

Have you ever needed to debug issues and ended up querying Prometheus for `group({app="foo"}) by (__name__)` to find out
what metrics the app exposes and than querying all of them fo find anything suspicious? Or do you often encounter apps
that do not have any official dashboard?

_Well I have a good news for you, Autograf have you covered!_ :tada:

## How does it work?
Autograf generates Grafana dashboard directly from `/metrics` or based on a metrics matching provided selector. Each
metric has own panel optimized for its type and those are grouped based on metric namespacing. If you want Autograf can
even upload the dashboard right your to a Grafana for you!

[autograf-2.webm](https://user-images.githubusercontent.com/6112562/178546235-7f9f815d-e843-4b0c-84dc-4fba2270eedc.webm)

## Installation
Using [prebuilt binaries](https://github.com/fusakla/autograf/releases/latest), [Docker
image](https://hub.docker.com/r/fusakla/autograf) of build it yourself.

```bash
go install github.com/fusakla/autograf@latest
```
or
```bash
make build
```

## Example
To see Autograf in action you can use the [demo example](./examples/demo/README.md). It is a simple docker-compose
that starts up Prometheus, Node exporter and Grafana. The Grafana instance is pre-configured with a datasource
pointing to the Proemtheus and service account. There is also an autograf.json config preset to test it with the setup.
See it's README for more details.

## How to use?

```bash
./autograf --help
Usage: autograf

Autograf generates Grafana dashboard from Prometheus metrics either read from a /metrics endpoint or queried from live Prometheus instance. The dashboard JSON is by default printed to stdout. But can also upload the dashboard directly to
your Grafana instance. You can configure most of the flags using config file. See the docs.

Example from /metrics:

    curl http://foo.bar/metrics | autograf --metrics-file -

Example from Prometheus query:

    GRAFANA_TOKEN=xxx autograf --prometheus-url http://prometheus.foo --selector '{app="foo"}' --grafana-url http://grafana.bar

Flags:
  -h, --help                                       Show context-sensitive help.
      --debug                                      Enable debug logging
      --version                                    Print Autograf version and exit
  -i, --ignore-config                              Ignore any config file
  -f, --metrics-file=STRING                        File containing the metrics exposed by app (will read stdin if se to - )
      --open-metrics-format                        Metrics data are in the application/openmetrics-text format.
  -p, --prometheus-url=STRING                      URL of Prometheus instance to fetch the metrics from.
      --prometheus-bearer-token=STRING             Bearer token to use for authentication with Prometheus instance.
  -s, --selector=STRING                            Selector to filter metrics from the Prometheus instance.
      --grafana-variables=GRAFANA-VARIABLES,...    Labels used as a variables for filtering in dashboard
      --grafana-url=STRING                         URL of Grafana to upload the dashboard to, if not specified, dashboard JSON is printed to stdout
      --grafana-folder=STRING                      Name of target Grafana folder
      --grafana-dashboard-name=STRING              Name of the Grafana dashboard
      --grafana-data-source=STRING                 Name of the Grafana datasource to use
      --open-browser                               Open the Grafana dashboard automatically in browser
```

### Loading data from all metrics exposed by app
```bash
curl -q http://demo.robustperception.io:9090/metrics | ./autograf --metrics-file -
```

### Loading data from live Prometheus instance
Print Grafana dashboard JSON for all metrics matching selector `{instance="demo.do.prometheus.io:9090",
job="prometheus"}` from the configured Prometheus instance.
```bash
autograf --prometheus-url http://demo.robustperception.io:9090 --selector '{instance="demo.do.prometheus.io:9090", job="prometheus"}'
```

### Uploading dashboard directly to Grafana
```bash
GRAFANA_TOKEN="xxx" autograf --prometheus-url http://demo.robustperception.io:9090 --selector '{instance="demo.do.prometheus.io:9090", job="prometheus"}' --grafana-url https://foo.bar --grafana-folder test
Dashboard successfully generated, see https://grafana.foo.bar/d/ygUo8se7k/autograf-dashboard
```

## Config file
If you do not want to set all the flags again and again you can use a config file. By default autograf looks for it in
`~/.autograf.json` and `~/.config/autograf.json` but can be changed using the `AUTOGRAF_CONFIG` env variable.
See the [example](./examples/demo/autograf.json) used in the demo.

### Config file syntax
```json
{
    "prometheus_url": "http://demo.robustperception.io:9090",
    "prometheus_bearer_token": "xxx",

    "grafana_url": "https://grafana.foo.bar",
    "grafana_dashboard_name": "Autograf",
    "grafana_folder": "FUSAKLAS garbage",
    "grafana_datasource": "Prometheus",
    "grafana_token": "xxx",
    
    "open_browser": true
}
```

Than you can use simply just this!
```bash
autograf -s {job='foo'}
```

## Panel config customization (EXPERIMENTAL)
This feature allows you to customize how the panel will look like using the metric HELP text.
To use it include in the and of the metric HELP string ` AUTOGRAF:{...}` where the supported JSON options
can be found in the [`PanelConfig`](./packages/model/panel_config.go#L10). Example of such metric HELP can
be found in the [`./examples/metrics_custom.txt`](./examples/metrics_custom.txt).


## Future ideas
- **Autograf should actually be Grafana app plugin and user could just go to `https://grafana.foo.bar/autograf?selector={foo="bar"}` and
  the dashboard would show up right in the Grafana itself. Unfortunately my JS juju is not good enough for this.**
- Add custom visuals for well known metrics.
