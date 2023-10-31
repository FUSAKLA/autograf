package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/fusakla/autograf/pkg/generator"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

var (
	specialSuffixRegexp = regexp.MustCompile(`(.+)_(total|info|sum|count|bucket)`)
)

func stripSpecialSuffixes(metricName string) string {
	return specialSuffixRegexp.ReplaceAllString(metricName, "$1")
}

func NewClient(logger logrus.FieldLogger, prometheusURL string, transport http.RoundTripper) (*Client, error) {
	cfg := api.Config{Address: prometheusURL}
	if transport != nil {
		cfg.RoundTripper = transport
	}
	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		logger: logger,
		v1api:  v1.NewAPI(cli),
	}, nil
}

type Client struct {
	v1api  v1.API
	logger logrus.FieldLogger
}

func (c *Client) queryMetricsMetadata(ctx context.Context) (map[string][]v1.Metadata, error) {
	c.logger.Info("querying prometheus metric metadata")
	res, err := c.v1api.Metadata(ctx, "", "")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) queryMetricsForSelector(ctx context.Context, selector string) ([]*model.Sample, error) {
	query := fmt.Sprintf("group(%s) by (__name__)", selector)
	c.logger.WithField("query", query).Info("querying prometheus")
	res, warnings, err := c.v1api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		c.logger.WithField("warnings", warnings).Warn("encountered warnings while querying Prometheus")
	}
	switch r := res.(type) {
	case model.Vector:
		return r, nil
	default:
		return nil, fmt.Errorf("unexpected result type %s expected vector", res.Type().String())
	}
}

func (c *Client) MetricsForSelector(ctx context.Context, selector string) (map[string]*generator.Metric, error) {
	samples, err := c.queryMetricsForSelector(ctx, selector)
	if err != nil {
		return nil, err
	}
	metadata, err := c.queryMetricsMetadata(ctx)
	if err != nil {
		return nil, err
	}
	metrics := map[string]*generator.Metric{}
	for _, s := range samples {
		metricName := string(s.Metric["__name__"])
		metricMetadata, ok := metadata[stripSpecialSuffixes(metricName)]
		if !ok {
			metricMetadata = metadata[metricName]
		}
		if len(metricMetadata) > 0 {
			metrics[metricName] = &generator.Metric{
				Name:       metricName,
				MetricType: generator.MetricType(metricMetadata[0].Type),
				Help:       metricMetadata[0].Help,
				Unit:       generator.MetricUnit(metricMetadata[0].Unit),
			}
		} else {
			metrics[metricName] = &generator.Metric{
				Name: metricName,
			}
		}
	}
	return metrics, nil
}
