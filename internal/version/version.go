//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package version

import "fmt"

var (
	// Version holds the current version from VERSION file
	Version = "dev"
	// Commit holds the git commit hash
	Commit = "none"
	// BuildTime holds the build timestamp
	BuildTime = "unknown"
)

// String returns a formatted version string
func String() string {
	return fmt.Sprintf("%s (commit: %s, built at: %s)", Version, Commit, BuildTime)
}

// Short returns just the version number
func Short() string {
	return Version
}
