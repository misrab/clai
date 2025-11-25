package version

import (
	"fmt"
	"runtime"
)

const (
	// Number is the semantic version of the CLI.
	Number = "0.1.0"
	// GitCommit can be overridden at build time: go build -ldflags "-X github.com/misrab/clai/internal/version.GitCommit=abc123"
	GitCommit = "dev"
)

// Full returns a printable, detailed version string.
func Full() string {
	return fmt.Sprintf("clai %s (%s) %s/%s", Number, GitCommit, runtime.GOOS, runtime.GOARCH)
}
