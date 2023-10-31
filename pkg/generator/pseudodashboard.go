package generator

import "encoding/json"

func NewPseudoDashboardFromMetrics(metrics map[string]*Metric) PseudoDashboard {
	return groupIntoPseudoDashboard(metrics)
}

type PseudoDashboard struct {
	Rows map[string]*PseudoRow `json:"rows"`
}

func (d PseudoDashboard) String() string {
	data, _ := json.MarshalIndent(d, "", "  ")
	return string(data)
}

func (d *PseudoDashboard) AddRowPanels(rowName string, metrics []*Metric) {
	if d.Rows == nil {
		d.Rows = make(map[string]*PseudoRow)
	}
	row, ok := d.Rows[rowName]
	if !ok {
		d.Rows[rowName] = &PseudoRow{Metrics: metrics}
	} else {
		row.Metrics = append(row.Metrics, metrics...)
	}
}

type PseudoRow struct {
	Metrics []*Metric `json:"panels"`
}
