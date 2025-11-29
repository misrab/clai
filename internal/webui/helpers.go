package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// getURLParam extracts a URL parameter from chi router
func getURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// APIError represents a structured error response
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// respondJSON writes a JSON response with proper headers
func respondJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// respondError writes a JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIError{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// decodeJSON decodes a JSON request body into the target
func decodeJSON(r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	return nil
}

// validateRequired validates that required fields in a struct are not empty
// Returns error message if validation fails
func validateRequired(fields map[string]string) error {
	var missing []string
	for name, value := range fields {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required fields missing: %s", strings.Join(missing, ", "))
	}
	return nil
}

// setupSSE configures response headers for Server-Sent Events
// func setupSSE(w http.ResponseWriter) {
// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// }

// // writeSSE writes data in SSE format
// func writeSSE(w http.ResponseWriter, data interface{}) error {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Fprintf(w, "data: %s\n\n", jsonData)
// 	return nil
// }
