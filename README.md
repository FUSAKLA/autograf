# Autograf
[![Go Report
Card](https://goreportcard.com/badge/github.com/fusakla/autograf)](https://goreportcard.com/report/github.com/fusakla/autograf)
[![GitHub actions
CI](https://img.shields.io/github/workflow/status/fusakla/autograf/Go/master)](https://github.com/FUSAKLA/autograf/actions?query=branch%3Amain)
[![Docker Pulls](https://img.shields.io/docker/pulls/fusakla/autograf)](https://hub.docker.com/r/fusakla/autograf)
[![GitHub binaries
download](https://img.shields.io/github/downloads/fusakla/autograf/total?label=Prebuilt%20binaries%20downloads)](https://github.com/FUSAKLA/autograf/releases/latest)

**Dynamically generate Grafana dashboard based on Prometheus metrics**

<p align="center"><img src="./autograf.excalidraw.png"></p>

Have you ever needed to debug issues and ended up querying Prometheus for `group({app="foo"}) by (__name__)` to find out
what metrics the app exposes and than querying all of them fo find anything suspicious? Or do you often encounter apps
that do not have any official dashboard?

_Well I have a good news for you, Autograf have you covered!_ :tada:

## How does it work?
Autograf generates Grafana dashboard directly form `/metrics` or based on a metrics matching provided selector. Each
metric has own panel optimized for its type and those are grouped based on metric namespacing. If you want Autograf can
event upload the dashboard Right your to a Grafana for you!

[autograf-2.webm](https://user-images.githubusercontent.com/6112562/178546235-7f9f815d-e843-4b0c-84dc-4fba2270eedc.webm)

## Installation
Using [prebuilt binaries](https://github.com/FUSAKLA/autograf/releases/latest), [Docker
image](https://hub.docker.com/r/fusakla/autograf) of build it yourself.

```bash
go install github.com/fusakla/autograf@latest
```
or
```bash
make build
```

## How to use?

```bash
./autograf --help
Usage: autograf

Autograf generates Grafana dashboard form Prometheus metrics either read form a /metrics endpoint or queried form live Prometheus instance. The dashboard JSON is by default printed to stdout. But can also upload
the dashboard directly to your Grafana instance. You can configure most of the flags using config file. See the docs.

Example from /metrics:

    curl http://foo.bar/metrics | autograf --metrics-file -

Example from Prometheus query:

    GRAFANA_TOKEN=xxx autograf --prometheus-url http://prometheus.foo --selector {app='foo'} --grafana-url http://grafana.bar

Flags:
  -h, --help                                       Show context-sensitive help.
      --debug                                      Enable debug logging
      --metrics-file=STRING                        File containing the metrics exposed by app (will read stdin if se to - )
      --open-metrics-format                        Metrics data are in the application/openmetrics-text format.
      --prometheus-url=STRING                      URL of Prometheus instance to fetch the metrics from.
      --selector=STRING                            Selector to filter metrics from the Prometheus instance.
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
job="prometheus"}` form the configured Prometheus instance.
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
`~/.autograf.json` but can be changed using the `AUTOGRAF_CONFIG` env variable.

### Config file syntax
```json
{
    "prometheus_url": "http://demo.robustperception.io:9090",
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
autograf --selector {app='foo'}
```


## Future ideas
- Could be a Grafana app plugin and user could just go to `https://grafana.foo.bar/autograf?selector={foo="bar"}` and
  the dashboard would show up right in the Grafana itself.
- Allow customizing the graph visualization using some tags in metric HELP(panel type, aggregations, units, colors,
  description, ...)
- Add custom visuals for well known metrics 
