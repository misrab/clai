package webui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/misrab/clai/internal/ai"
	"github.com/misrab/clai/internal/storage"
)

// HandleListChats handles GET /api/chats
func HandleListChats(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chats, err := store.ListChats()
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list chats: %v", err))
			return
		}

		respondJSON(w, http.StatusOK, chats)
	}
}

// HandleCreateChat handles POST /api/chats
func HandleCreateChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}

		if err := decodeJSON(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := validateRequired(map[string]string{"id": req.ID}); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if req.Title == "" {
			req.Title = "New Chat"
		}

		now := time.Now()
		chat := &storage.Chat{
			ID:        req.ID,
			Title:     req.Title,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := store.CreateChat(chat); err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create chat: %v", err))
			return
		}

		respondJSON(w, http.StatusCreated, chat)
	}
}

// HandleGetChat handles GET /api/chats/{id}
func HandleGetChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getURLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "Chat ID is required")
			return
		}

		chat, messages, err := store.GetChatWithMessages(id)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get chat: %v", err))
			return
		}

		if chat == nil {
			respondError(w, http.StatusNotFound, "Chat not found")
			return
		}

		response := struct {
			*storage.Chat
			Messages []*storage.Message `json:"messages"`
		}{
			Chat:     chat,
			Messages: messages,
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// HandleUpdateChat handles PUT /api/chats/{id}
func HandleUpdateChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getURLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "Chat ID is required")
			return
		}

		var req struct {
			Title string `json:"title"`
		}

		if err := decodeJSON(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := validateRequired(map[string]string{"title": req.Title}); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := store.UpdateChatTitle(id, req.Title); err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update chat: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// HandleDeleteChat handles DELETE /api/chats/{id}
func HandleDeleteChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getURLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "Chat ID is required")
			return
		}

		if err := store.DeleteChat(id); err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete chat: %v", err))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// sendMessageRequest represents the request body for sending a message
type sendMessageRequest struct {
	UserMessageID string `json:"userMessageId"`
	Content       string `json:"content"`
	Model         string `json:"model,omitempty"`
}

// HandleSendMessage handles POST /api/chats/{id}/send
// Saves user message, gets AI response from Ollama, and streams or returns the response
func HandleSendMessage(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := getURLParam(r, "id")
		if chatID == "" {
			respondError(w, http.StatusBadRequest, "Chat ID is required")
			return
		}

		var req sendMessageRequest
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := validateRequired(map[string]string{
			"userMessageId": req.UserMessageID,
			"content":       req.Content,
		}); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Save user message
		userMessage := &storage.Message{
			ID:        req.UserMessageID,
			ChatID:    chatID,
			Role:      "user",
			Content:   req.Content,
			CreatedAt: time.Now(),
		}

		if err := store.CreateMessage(userMessage); err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save user message: %v", err))
			return
		}

		// Use specified model or default
		model := req.Model
		if model == "" {
			model = "codellama:7b"
		}

		// Server decides whether to stream (default: always stream for now)
		if shouldStream(req.Content) {
			handleStreamingResponse(w, r, chatID, req.Content, model, store)
		} else {
			handleNonStreamingResponse(w, chatID, req.Content, model, store)
		}
	}
}

// shouldStream determines if the response should be streamed
// For now, always returns true. In the future, can check content type, metadata, etc.
func shouldStream(content string) bool {
	// Future logic:
	// - Check if response will contain HTML/attachments -> return false
	// - Check if it's a code generation request -> return true
	// - Check message metadata/flags -> return accordingly
	return true
}

// handleStreamingResponse streams the AI response using SSE
func handleStreamingResponse(w http.ResponseWriter, r *http.Request,
	chatID, content, model string, store *storage.Store) {

	setupSSE(w)

	flusher, ok := w.(http.Flusher)
	if !ok {
		respondError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	client := ai.NewClient(model)
	assistantID := generateMessageID()
	var fullResponse string

	// Stream chunks to client and accumulate full response
	err := client.ChatStream(content, func(chunk string) error {
		// Check if client disconnected
		if r.Context().Err() != nil {
			return fmt.Errorf("client disconnected")
		}

		fullResponse += chunk

		// Send chunk via SSE
		if err := writeSSE(w, map[string]interface{}{
			"id":    assistantID,
			"chunk": chunk,
			"done":  false,
		}); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	})

	if err != nil {
		writeSSE(w, map[string]interface{}{
			"error": err.Error(),
			"done":  true,
		})
		flusher.Flush()
		return
	}

	// Save complete message to DB
	assistantMessage := &storage.Message{
		ID:        assistantID,
		ChatID:    chatID,
		Role:      "assistant",
		Content:   fullResponse,
		CreatedAt: time.Now(),
	}

	if err := store.CreateMessage(assistantMessage); err != nil {
		fmt.Printf("Failed to save assistant message: %v\n", err)
	}

	// Send final event with full message
	writeSSE(w, map[string]interface{}{
		"id":      assistantID,
		"content": fullResponse,
		"done":    true,
	})
	flusher.Flush()
}

// handleNonStreamingResponse returns the complete AI response at once
func handleNonStreamingResponse(w http.ResponseWriter,
	chatID, content, model string, store *storage.Store) {

	client := ai.NewClient(model)
	aiResponse, err := client.Chat(content)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Save assistant message
	assistantMessage := &storage.Message{
		ID:        generateMessageID(),
		ChatID:    chatID,
		Role:      "assistant",
		Content:   aiResponse,
		CreatedAt: time.Now(),
	}

	if err := store.CreateMessage(assistantMessage); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save assistant message: %v", err))
		return
	}

	respondJSON(w, http.StatusCreated, assistantMessage)
}

// generateMessageID generates a random message ID
func generateMessageID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
