//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package logger

import (
	"fmt"

	"github.com/dcvix/dcvix-stats/internal/globals"
)

func LogVerbose(format string, args ...interface{}) {
	if globals.Verbose {
		fmt.Printf(" "+format, args...)
	}
}
