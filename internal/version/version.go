package version

import (
	"fmt"
	"runtime"
)

// These variables are set via ldflags during build time by GoReleaser
var (
	// Version is the semantic version of the CLI (e.g., "v0.0.4")
	Version = "0.0.0-dev"
	// Commit is the git commit hash
	Commit = "dev"
	// Date is the build date
	Date = "unknown"
)

// Full returns a printable, detailed version string.
func Full() string {
	return fmt.Sprintf("clai %s (commit: %s, built: %s) %s/%s", Version, Commit, Date, runtime.GOOS, runtime.GOARCH)
}
