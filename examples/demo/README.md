# Example demo of Autograf

Yo will need to install or build the autograf, see [the installation instructions](../../README.md#installation).


Than you can just run
```bash
# Spin-up Prometheus and Grafana
docker-compose up -d
# Wait for it to be ready
sleep 30
# Generate your first dashboard using Autograf!
AUTOGRAF_CONFIG=autograf.json autograf -s '{job="node-exporter"}'
```

To cleanup just run
```
docker-compose rm -f -s -v
```
