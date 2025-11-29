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

// HandleSendMessage handles POST /api/chats/{id}/send
// This endpoint saves the user message, gets AI response from Ollama, and saves the assistant response
func HandleSendMessage(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := getURLParam(r, "id")
		if chatID == "" {
			respondError(w, http.StatusBadRequest, "Chat ID is required")
			return
		}

		var req struct {
			UserMessageID string `json:"userMessageId"`
			Content       string `json:"content"`
			Model         string `json:"model,omitempty"`
		}

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

		// Get AI response from Ollama (non-streaming for simplicity)
		client := ai.NewClient(model)
		aiResponse, err := client.Chat(req.Content)
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
}

// generateMessageID generates a random message ID
func generateMessageID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
