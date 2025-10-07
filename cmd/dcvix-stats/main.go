//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2/app"

	"github.com/dcvix/dcvix-stats/internal/globals"
	"github.com/dcvix/dcvix-stats/internal/gui"
	"github.com/dcvix/dcvix-stats/internal/logger"
	"github.com/dcvix/dcvix-stats/internal/version"
)

// var WinWidth = 800
// var WinHeigh = 600
var WinWidth = 512
var WinHeigh = 384

func main() {

	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(&globals.Verbose, "verbose", false, "Enable verbose logging")
	flag.IntVar(&globals.LogEntriesQty, "entries", 120, "How many last entries/minutes to evaluate")
	flag.StringVar(&globals.LogFile, "logfile", getDefaultLogPath(), "Path to the DCV server log file")
	flag.IntVar(&globals.RefreshInterval, "refresh", 30, "Auto-refresh interval in seconds")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\n", version.String())
		os.Exit(0)
	}

	logger.LogVerbose("Starting log parser for file: %s\n", globals.LogFile)
	logger.LogVerbose("Refreshing every %v seconds...\n", globals.RefreshInterval)

	// setup main window.
	a := app.New()
	w := gui.NewMainWindow(a)
	w.ShowAndRun()
}

func getDefaultLogPath() string {
	if runtime.GOOS == "windows" {
		return `C:\ProgramData\NICE\dcv\log\server.log`
	}
	return "/var/log/dcv/server.log"
}
