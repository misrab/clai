package storage

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDBPath returns the path to the SQLite database following OS conventions
func GetDBPath() (string, error) {
	var dataDir string

	switch runtime.GOOS {
	case "windows":
		dataDir = os.Getenv("LOCALAPPDATA")
		if dataDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataDir = filepath.Join(home, "AppData", "Local")
		}
	default: // linux, darwin
		dataDir = os.Getenv("XDG_DATA_HOME")
		if dataDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataDir = filepath.Join(home, ".local", "share")
		}
	}

	appDir := filepath.Join(dataDir, "clai")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "clai.db"), nil
}

