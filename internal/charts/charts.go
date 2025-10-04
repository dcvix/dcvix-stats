// SPDX-FileCopyrightText: 2025 Diego Cortassa
// SPDX-License-Identifier: MIT

package charts

import (
	"fmt"

	"github.com/vicanso/go-charts/v2"

	"github.com/dcvix/dcvix-stats/internal/logparser"
)

func ChartByMetricList(metrics map[string][]logparser.LogEntry, width float32, height float32) []byte {
	var values [][]float64
	var timeStamps []string

	for _, entries := range metrics {
		var metricValues []float64
		var metricTimeStamps []string

		for _, entry := range entries {
			metricValues = append(metricValues, entry.LastValue)
			metricTimeStamps = append(metricTimeStamps, entry.Timestamp)
		}

		values = append(values, metricValues)
		timeStamps = append(timeStamps, metricTimeStamps...)
	}

	return Chart([]string{""}, values, timeStamps, width, height)
}

func Chart(metrics []string, values [][]float64, timeStamps []string, width float32, height float32) []byte {
	p, err := charts.LineRender(
		values,
		// charts.TitleTextOptionFunc("Line"),
		charts.XAxisDataOptionFunc(timeStamps),
		charts.LegendLabelsOptionFunc(metrics, "100"),
		func(opt *charts.ChartOption) {
			opt.Theme = "grafana"
			opt.Legend.Padding = charts.Box{
				Top:    5,
				Bottom: 10,
			}
			opt.YAxisOptions = []charts.YAxisOption{
				{
					SplitLineShow: charts.FalseFlag(),
				},
			}
			opt.SymbolShow = charts.FalseFlag()
			opt.LineStrokeWidth = 1
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
			opt.Width = int(width)
			opt.Height = int(height)
		},
	)

	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	return buf
}
