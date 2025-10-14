//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package logger

import (
	"log"

	"github.com/dcvix/dcvix-stats/internal/globals"
)

func LogVerbose(v ...interface{}) {
	if globals.Verbose {
		log.Println(v...)
	}
}

func LogVerbosef(format string, v ...interface{}) {
	if globals.Verbose {
		log.Printf(format, v...)
	}
}
