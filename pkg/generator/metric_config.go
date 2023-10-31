package generator

import (
	"encoding/json"
	"strings"

	validate "github.com/go-playground/validator/v10"
)

type MetricConfig struct {
	// Overrides name of the row, all metrics with the same row name will be grouped together
	Row string

	// Type of aggregation to use
	Aggregation string `validate:"omitempty,oneof=avg max min group count sum"`
	// Labels to aggregate by
	AggregateBy []string `json:"aggregate_by" validate:"excluded_without=Aggregation"`

	// Stack series
	Stack bool
	// How thick the line should be
	LineWidth int `validate:"gte=0,lte=10"`
	// Opacity of line fill
	Fill int `validate:"gte=0,lte=100"`
	// Scale of the Y axis
	Scale string `validate:"oneof=linear log2 log10"`
	// Where to place the legend
	LegendPosition string `json:"legend_position" validate:"oneof=bottom right hide"`
	// What calculations include in the legend
	LegendCalcs []string `json:"legend_calculations"`
	// Width of the panel
	Width int `validate:"gte=1,lte=12"`
	// Height of the panel
	Height int `validate:"gte=1,lte=12"`
	// Metric name to use as a max threshold in the graph
	MaxFromMetric string `json:"max_from_metric"`
	// Metric name to use as a min threshold in the graph
	MinFromMetric string `json:"min_from_metric"`
}

func (c MetricConfig) ThresholdMetric() *Metric {
	defaultConfig := MetricConfig{LineWidth: 2, Fill: 0, Stack: false}
	if c.MaxFromMetric != "" {
		defaultConfig.Aggregation = "min"
		return &Metric{Name: c.MaxFromMetric, Config: defaultConfig}
	} else if c.MinFromMetric != "" {
		defaultConfig.Aggregation = "max"
		return &Metric{Name: c.MaxFromMetric, Config: defaultConfig}
	}
	return nil
}

var configSeparator = " AUTOGRAF:"

var defaultConfig = MetricConfig{
	Stack:          false,
	LineWidth:      1,
	Fill:           1,
	Scale:          "linear",
	LegendPosition: "bottom",
	LegendCalcs:    []string{"max", "avg", "last"},
	Width:          4,
	Height:         5,
}

func LoadConfigFromHelp(help string) (MetricConfig, error) {
	conf := defaultConfig
	parts := strings.Split(help, configSeparator)
	if len(parts) < 2 {
		return conf, nil
	}
	dec := json.NewDecoder(strings.NewReader(strings.TrimSpace(parts[1])))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&conf); err != nil {
		return conf, err
	}
	err := validate.New().Struct(conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
