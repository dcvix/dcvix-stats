//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

//go:generate fyne bundle --package gui -o bundled.go  ../../assets/icon.png

package gui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/dcvix/dcvix-stats/internal/globals"
	"github.com/dcvix/dcvix-stats/internal/logger"
	"github.com/dcvix/dcvix-stats/internal/logparser"
	"github.com/dcvix/dcvix-stats/internal/version"
)

const WinWidth = 512
const WinHeigh = 384

func NewMainWindow(a fyne.App) fyne.Window {

	parser := logparser.NewLogParser(globals.LogFile)

	w := a.NewWindow(globals.AppName)

	// Set App Icon
	w.SetIcon(resourceIconPng)

	// ## System Tray
	// if desk, ok := a.(desktop.App); ok {
	// 	menu := fyne.NewMenu(globals.AppName,
	// 		fyne.NewMenuItem("Show", func() {
	// 			log.Println("show")
	// 			w.Show()
	// 		}),
	// 		fyne.NewMenuItem("Quit", func() {
	// 			a.Quit()
	// 		}),
	// 	)
	// 	desk.SetSystemTrayMenu(menu)
	// }
	// // When using tray menu we only hide the main window when closed
	// w.SetCloseIntercept(func() { w.Hide() })

	// ## Main menu
	parseURL := func(urlStr string) *url.URL {
		link, err := url.Parse(urlStr)
		if err != nil {
			fyne.LogError("Could not parse URL", err)
		}
		return link
	}

	showAbout := func() {
		aboutWindow := a.NewWindow("About")
		aboutWindow.SetContent(container.NewVBox(
			widget.NewLabelWithStyle(globals.AppName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fmt.Sprintf("Version: %s", version.String())),
			widget.NewLabel("Author: Diego Cortassa"),
			widget.NewLabel("shows DCV connection statistics."),
			container.NewHBox(
				widget.NewLabel("License:"),
				widget.NewHyperlink("MIT", parseURL("https://github.com/dcvix/dcvix-stats/blob/main/LICENSE.md")),
			),
			container.NewHBox(
				widget.NewLabel("More info:"),
				widget.NewHyperlink("dcvix-stats", parseURL("https://github.com/dcvix/dcvix-stats")),
			),
		))
		aboutWindow.Resize(fyne.NewSize(300, 200))
		aboutWindow.Show()
	}

	var mainMenu *fyne.MainMenu

	// ## Main window setup
	w.Resize(fyne.Size{Width: float32(WinWidth), Height: float32(WinHeigh)})

	parser.ReadLogFile()

	// res := parser.GetEntriesByMetric("quic_lost_packets")
	// for _, entry := range res {
	// 	fmt.Printf("*** %v %v\n", entry.Timestamp, entry.LastValue)
	// }

	values, timeStamps := parser.GetLatestData()
	logger.LogVerbose("Timestamps (timeStamps): %v\n", timeStamps)
	logger.LogVerbose("Timestamps (values): %v\n", values)

	for i, metric := range globals.Metrics {
		if i < len(values) && len(values[i]) > 0 {
			logger.LogVerbose("%-25s(metric): %v\n", metric, values[i])
		}
	}

	type graphConfig struct {
		name             string
		metrics          []string
		chartView        *ChartView
		menuItem         *fyne.MenuItem
		enabledByDefault bool
	}

	graphConfigs := []*graphConfig{
		{
			name:             "QUICLostPktsGraph",
			metrics:          []string{"quic_lost_packets", "quic_lost_packets_avg"},
			enabledByDefault: true,
		},
		{
			name:             "QUICSentRecvPktsGraph",
			metrics:          []string{"quic_sent_packets_avg", "quic_recv_packets_avg"},
			enabledByDefault: true,
		},
		{
			name:             "QUICRttNanos",
			metrics:          []string{"quic_rtt_nanos", "quic_rtt_nanos_avg"},
			enabledByDefault: true,
		},
		{
			name:             "QUICCwndSize",
			metrics:          []string{"quic_cwnd_size", "quic_cwnd_size_avg"},
			enabledByDefault: true,
		},
		{
			name:             "QUICDeliveryRate",
			metrics:          []string{"quic_delivery_rate", "quic_delivery_rate_avg"},
			enabledByDefault: false,
		},
		{
			name:             "DGrams",
			metrics:          []string{"dgram_sent", "dgram_sent_avg", "dgram_recv", "dgram_recv_avg"},
			enabledByDefault: false,
		},
		{
			name:             "StreamsGraph",
			metrics:          []string{"active_streams", "stream_sent", "stream_recv"},
			enabledByDefault: false,
		},
		{
			name:             "ActiveStreamsGraph",
			metrics:          []string{"active_streams"},
			enabledByDefault: false,
		},
	}

	showMenuItems := make([]*fyne.MenuItem, 0, len(graphConfigs))
	graphContainers := make([]fyne.CanvasObject, 0, len(graphConfigs))

	for _, config := range graphConfigs {
		values, timeStamps := parser.GetEntriesByMetricList(config.metrics)
		config.chartView = NewChartView(config.metrics, values, timeStamps)
		graphContainers = append(graphContainers, config.chartView)

		config.menuItem = fyne.NewMenuItem(config.name, func() {
			config.menuItem.Checked = !config.menuItem.Checked
			if config.menuItem.Checked {
				config.chartView.Show()
			} else {
				config.chartView.Hide()
			}
			w.SetMainMenu(mainMenu)
		})
		config.menuItem.Checked = config.enabledByDefault
		if !config.menuItem.Checked {
			config.chartView.Hide()
		}
		showMenuItems = append(showMenuItems, config.menuItem)
	}

	refresh := func() {
		parser.ReadLogFile()
		for _, config := range graphConfigs {
			values, timeStamps := parser.GetEntriesByMetricList(config.metrics)
			config.chartView.RefreshData(values, timeStamps)
		}
	}

	// Menus
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("About", showAbout),
		fyne.NewMenuItem("Refresh", refresh),
	)

	showMenu := fyne.NewMenu("Show", showMenuItems...)
	mainMenu = fyne.NewMainMenu(fileMenu, showMenu)
	w.SetMainMenu(mainMenu)

	// Main container
	w.SetContent(container.NewAdaptiveGrid(2, graphContainers...))

	// go parser.Run()
	// parser.Run()
	return w

}

// ticker := time.NewTicker(30 * time.Second)
// defer ticker.Stop()

// for {
// <-ticker.C
// select {
// case <-ticker.C:
// 	// Continue to next iteration
// }
