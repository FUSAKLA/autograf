# Autograf
**Dynamically generate Grafana dashboard based on Prometheus metrics**


![](./autograf.excalidraw.png)

Have you ever needed to debug issues and ended up querying Prometheus for `group({app="foo"}) by (__name__)` to find out
what metrics the app exposes and than querying all of them fo find anything suspicious? Or do you often encounter apps
that do not have any official dashboard?

_Well I have a good news for you, Autograf have you covered!_ :tada:

## How does it work?
Autograf looks for all the metrics exposed by app on the `/metrics` endpoint or loads them from Prometheus and than
based on its metadata generated Grafana dashboard for you with panel for each metric grouped into rows based on their
namespacing. If you want Autograf can event upload the dashboard Right your to a Grafana for you!

## Installation
```
go install github.com/fusakla/autograf@latest
```

## How to use?

```bash
./autograf generate --help
Usage: autograf generate --metrics-file=METRICS-FILE --prometheus-url=PROMETHEUS-URL

Generates Grafana dashboard based on a given Prometheus metrics and prints it to stdout if not specified otherwise.

Flags:
  -h, --help                                   Show context-sensitive help.
      --debug                                  Enable debug logging.

      --open-metrics-format                    Metrics data are in the application/openmetrics-text format.
      --selector="{}"                          Selector to filter metrics from the Prometheus instance.
      --aggregate-by=AGGREGATE-BY,...          Additional labels to aggregate the queries by.
      --dashboard-name="Autograf dashboard"    Name of the dashboard
      --data-source=STRING                     Name of the data source to use
      --grafana-url=GRAFANA-URL                URL of Grafana to upload the dashboard to, if not specified, dashboard JSON is printed to stdout
      --grafana-folder=STRING                  Name of target Grafana folder

input
  --metrics-file=METRICS-FILE        File containing the metrics exposed by app (will read stdin if se to -)
  --prometheus-url=PROMETHEUS-URL    URL of Prometheus instance to fetch the metrics from.
```

### Loading data from all metrics exposed by app
```bash
curl -q http://demo.robustperception.io:9090/metrics | ./autograf generate --metrics-file -
```

### Loading data from live Prometheus instance
Print Grafana dashboard JSON for all metrics matching selector `{instance="demo.do.prometheus.io:9090",
job="prometheus"}` form the configured Prometheus instance.
```bash
autograf generate --prometheus-url http://demo.robustperception.io:9090 --selector '{instance="demo.do.prometheus.io:9090", job="prometheus"}'
```

### Uploading directly to Grafana
```bash
GRAFANA_TOKEN="xxx" autograf generate --prometheus-url http://demo.robustperception.io:9090 --selector '{instance="demo.do.prometheus.io:9090", job="prometheus"}' --grafana-url https://foo.bar --grafana-folder test
Dashboard successfully generated, see https://grafana.foo.bar/d/ygUo8se7k/autograf-dashboard
```

## Configuration
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

## Future ideas
- Could be a Grafana app plugin and user could just go to `https://grafana.foo.bar/autograf?selector={foo="bar"}` and
  the dashboard would show up right in the Grafana itself.
- Allow customizing the graph visualization using some tags in metric HELP(panel type, aggregations, units, colors,
  description, ...)
- Add custom visuals for well known metrics 
