package webui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/misrab/clai/internal/storage"
)

// Start starts the web UI server with the provided embedded filesystem
func Start(distFiles embed.FS, port int, openBrowser bool) error {
	// Initialize storage
	store, err := storage.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	// Get the dist subdirectory from embedded files
	distFS, err := fs.Sub(distFiles, "webui/dist")
	if err != nil {
		return fmt.Errorf("failed to access embedded files: %w", err)
	}

	// Create chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	}))

	// API routes
	r.Route("/api/chats", func(r chi.Router) {
		r.Get("/", HandleListChats(store))
		r.Post("/", HandleCreateChat(store))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", HandleGetChat(store))
			r.Put("/", HandleUpdateChat(store))
			r.Delete("/", HandleDeleteChat(store))
			r.Post("/send", HandleSendMessage(store))
		})
	})

	// Serve embedded files
	r.Handle("/*", http.FileServer(http.FS(distFS)))

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
	return http.ListenAndServe(":"+fmt.Sprintf("%d", port), r)
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
