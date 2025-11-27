package webui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// Start starts the web UI server with the provided embedded filesystem
func Start(distFiles embed.FS, port int, openBrowser bool) error {
	// Get the dist subdirectory from embedded files
	distFS, err := fs.Sub(distFiles, "webui/dist")
	if err != nil {
		return fmt.Errorf("failed to access embedded files: %w", err)
	}

	// Serve embedded files
	http.Handle("/", http.FileServer(http.FS(distFS)))

	addr := fmt.Sprintf("localhost:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	fmt.Printf("Starting clai web UI at %s\n", url)

	// Auto-open browser
	if openBrowser {
		// Wait a moment for server to start
		time.AfterFunc(500*time.Millisecond, func() {
			openURL(url)
		})
	}

	fmt.Println("Press Ctrl+C to stop")
	return http.ListenAndServe(":"+fmt.Sprintf("%d", port), nil)
}

// openURL opens a URL in the default browser
func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		fmt.Println("Please open:", url)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\nPlease open: %s\n", err, url)
	}
}

