package generator

import (
	"fmt"
	"strings"
)

type LimitType string

const (
	LimitMax LimitType = "max"
	LimitMin LimitType = "min"
)

func metricWithSelector(metric, selector string) string {
	return fmt.Sprintf("%s%s", metric, selector)
}

func aggregateQuery(query, aggregation string, aggragateBy []string) string {
	return fmt.Sprintf("%s(%s) by (%s)", aggregation, query, strings.Join(aggragateBy, ","))
}

func rateCounterQuery(query, rangeSelector string) string {
	return fmt.Sprintf("rate(%s[%s])", query, rangeSelector)
}

func ThresholdQuery(metricName string, selector string, lType LimitType) string {
	aggregation := "min"
	if lType == LimitMin {
		aggregation = "max"
	}
	return aggregateQuery(metricWithSelector(metricName, selector), aggregation, []string{})
}
