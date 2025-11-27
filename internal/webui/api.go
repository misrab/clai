package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/misrab/clai/internal/storage"
)

// HandleListChats handles GET /api/chats
func HandleListChats(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		chats, err := store.ListChats()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list chats: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chats)
	}
}

// HandleCreateChat handles POST /api/chats
func HandleCreateChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
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
			http.Error(w, fmt.Sprintf("Failed to create chat: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(chat)
	}
}

// HandleGetChat handles GET /api/chats/{id}
func HandleGetChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path
		id := strings.TrimPrefix(r.URL.Path, "/api/chats/")
		if id == "" {
			http.Error(w, "Chat ID is required", http.StatusBadRequest)
			return
		}

		chat, messages, err := store.GetChatWithMessages(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get chat: %v", err), http.StatusInternalServerError)
			return
		}

		if chat == nil {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}

		response := struct {
			*storage.Chat
			Messages []*storage.Message `json:"messages"`
		}{
			Chat:     chat,
			Messages: messages,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// HandleUpdateChat handles PUT /api/chats/{id}
func HandleUpdateChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path
		id := strings.TrimPrefix(r.URL.Path, "/api/chats/")
		if id == "" {
			http.Error(w, "Chat ID is required", http.StatusBadRequest)
			return
		}

		var req struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		if err := store.UpdateChatTitle(id, req.Title); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update chat: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// HandleDeleteChat handles DELETE /api/chats/{id}
func HandleDeleteChat(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path
		id := strings.TrimPrefix(r.URL.Path, "/api/chats/")
		if id == "" {
			http.Error(w, "Chat ID is required", http.StatusBadRequest)
			return
		}

		if err := store.DeleteChat(id); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete chat: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// HandleCreateMessage handles POST /api/chats/{id}/messages
func HandleCreateMessage(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract chat ID from path
		path := strings.TrimPrefix(r.URL.Path, "/api/chats/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 || parts[1] != "messages" {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		chatID := parts[0]

		var req struct {
			ID      string `json:"id"`
			Role    string `json:"role"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.ID == "" || req.Role == "" || req.Content == "" {
			http.Error(w, "ID, role, and content are required", http.StatusBadRequest)
			return
		}

		if req.Role != "user" && req.Role != "assistant" {
			http.Error(w, "Role must be 'user' or 'assistant'", http.StatusBadRequest)
			return
		}

		message := &storage.Message{
			ID:        req.ID,
			ChatID:    chatID,
			Role:      req.Role,
			Content:   req.Content,
			CreatedAt: time.Now(),
		}

		if err := store.CreateMessage(message); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create message: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}

// chatHandler routes requests based on path and method
func chatHandler(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/chats")
		
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Route based on path
		if path == "" || path == "/" {
			// /api/chats
			if r.Method == http.MethodGet {
				HandleListChats(store)(w, r)
			} else if r.Method == http.MethodPost {
				HandleCreateChat(store)(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if strings.Contains(path, "/messages") {
			// /api/chats/{id}/messages
			HandleCreateMessage(store)(w, r)
		} else {
			// /api/chats/{id}
			if r.Method == http.MethodGet {
				HandleGetChat(store)(w, r)
			} else if r.Method == http.MethodPut {
				HandleUpdateChat(store)(w, r)
			} else if r.Method == http.MethodDelete {
				HandleDeleteChat(store)(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}
}

