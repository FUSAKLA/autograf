package generator

import "testing"

func Test_metricWithSelector(t *testing.T) {
	type args struct {
		metric   string
		selector string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty selector",
			args: args{
				metric:   "queue_size",
				selector: "",
			},
			want: "queue_size",
		},
		{
			name: "Valid selector",
			args: args{
				metric:   "queue_size",
				selector: `{foo="bar"}`,
			},
			want: `queue_size{foo="bar"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metricWithSelector(tt.args.metric, tt.args.selector); got != tt.want {
				t.Errorf("metricWithSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_aggregateQuery(t *testing.T) {
	type args struct {
		query       string
		aggregation string
		aggragateBy []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty aggragateBy",
			args: args{
				query:       "queue_size",
				aggregation: "sum",
				aggragateBy: []string{},
			},
			want: "sum(queue_size) by ()",
		},
		{
			name: "Valid aggragateBy",
			args: args{
				query:       "queue_size",
				aggregation: "avg",
				aggragateBy: []string{"foo", "bar"},
			},
			want: "avg(queue_size) by (foo,bar)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aggregateQuery(tt.args.query, tt.args.aggregation, tt.args.aggragateBy); got != tt.want {
				t.Errorf("aggregateQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_rateCounterQuery(t *testing.T) {
	type args struct {
		query         string
		rangeSelector string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty rangeSelector",
			args: args{
				query:         "events_total",
				rangeSelector: "$__range_interval",
			},
			want: "rate(events_total[$__range_interval])",
		},
		{
			name: "Valid rangeSelector",
			args: args{
				query:         "events_total",
				rangeSelector: "5m",
			},
			want: "rate(events_total[5m])",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rateCounterQuery(tt.args.query, tt.args.rangeSelector); got != tt.want {
				t.Errorf("rateCounterQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
