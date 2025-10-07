//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

//go:generate fyne bundle --package gui -o bundled.go  ../../assets/icon.png

package gui

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/dcvix/dcvix-stats/internal/globals"
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
		config.chartView = NewChartView(config.metrics, nil, nil)
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

	// Reload log file and redraw graphs
	refresh := func() {
		err := parser.ReadLogFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error, could not read log file: %s\n", globals.LogFile)
			os.Exit(1)
		}

		for _, config := range graphConfigs {
			values, timeStamps := parser.GetEntriesByMetricList(config.metrics)
			config.chartView.RefreshData(values, timeStamps)
		}
	}
	refresh()

	// Auto refresh ticker
	var autoRefreshTicker *time.Ticker
	var autoRefreshDone chan bool

	startAutoRefresh := func() {
		if autoRefreshTicker != nil {
			return
		}
		autoRefreshTicker = time.NewTicker(time.Duration(globals.RefreshInterval) * time.Second)
		autoRefreshDone = make(chan bool)

		go func() {
			for {
				select {
				case <-autoRefreshTicker.C:
					refresh()
				case <-autoRefreshDone:
					return
				}
			}
		}()
	}

	stopAutoRefresh := func() {
		if autoRefreshTicker != nil {
			autoRefreshTicker.Stop()
			autoRefreshDone <- true
			close(autoRefreshDone)
			autoRefreshTicker = nil
		}
	}

	autoRefreshItem := fyne.NewMenuItem("Auto Refresh", nil)
	autoRefreshItem.Checked = false
	autoRefreshItem.Action = func() {
		autoRefreshItem.Checked = !autoRefreshItem.Checked
		if autoRefreshItem.Checked {
			startAutoRefresh()
		} else {
			stopAutoRefresh()
		}
		w.SetMainMenu(mainMenu)
	}

	// Window close handling
	w.SetCloseIntercept(func() {
		stopAutoRefresh()
		w.Close()
	})

	// Menus
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("About", showAbout),
		fyne.NewMenuItem("Refresh", refresh),
		autoRefreshItem,
	)

	showMenu := fyne.NewMenu("Show", showMenuItems...)
	mainMenu = fyne.NewMainMenu(fileMenu, showMenu)
	w.SetMainMenu(mainMenu)

	// Main container
	w.SetContent(container.NewAdaptiveGrid(2, graphContainers...))

	return w

}
