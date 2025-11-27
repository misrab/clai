package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/misrab/clai/internal/ai"
)

const defaultModel = "codellama:7b"

// chatRequest represents an incoming chat message
type chatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

// chatChunk represents a streaming response chunk
type chatChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

// HandleChat handles streaming chat requests using Ollama
func HandleChat(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	// Use specified model or default
	model := req.Model
	if model == "" {
		model = defaultModel
	}

	// Set headers for Server-Sent Events (SSE)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Create Ollama client (same as CLI)
	client := ai.NewClient(model)

	// Stream response from Ollama
	startTime := time.Now()
	fmt.Printf("[%s] Chat request: %q (model: %s)\n", startTime.Format("15:04:05"), req.Message, model)

	err := client.ChatStream(req.Message, func(chunk string) error {
		// Check if client disconnected
		if r.Context().Err() != nil {
			return fmt.Errorf("client disconnected")
		}

		// Send chunk as SSE
		data := chatChunk{
			Content: chunk,
			Done:    false,
		}

		if err := writeSSEChunk(w, data); err != nil {
			return err
		}
		flusher.Flush()

		return nil
	})

	// Handle errors
	if err != nil {
		fmt.Printf("[%s] Error: %v\n", time.Now().Format("15:04:05"), err)
		errorChunk := chatChunk{
			Content: "",
			Done:    true,
			Error:   err.Error(),
		}
		writeSSEChunk(w, errorChunk)
		flusher.Flush()
		return
	}

	// Send final "done" chunk
	finalChunk := chatChunk{
		Content: "",
		Done:    true,
	}
	writeSSEChunk(w, finalChunk)
	flusher.Flush()

	elapsed := time.Since(startTime)
	fmt.Printf("[%s] Chat completed in %v\n", time.Now().Format("15:04:05"), elapsed)
}

// writeSSEChunk writes a chunk in SSE format
func writeSSEChunk(w http.ResponseWriter, chunk chatChunk) error {
	data, err := json.Marshal(chunk)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
	return nil
}
