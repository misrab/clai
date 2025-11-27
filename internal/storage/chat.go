package storage

import (
	"database/sql"
	"time"
)

// Chat represents a chat conversation
type Chat struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateChat creates a new chat in the database
func (s *Store) CreateChat(chat *Chat) error {
	_, err := s.db.Exec(`
		INSERT INTO chats (id, title, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, chat.ID, chat.Title, chat.CreatedAt, chat.UpdatedAt)
	return err
}

// GetChat retrieves a chat by ID
func (s *Store) GetChat(id string) (*Chat, error) {
	chat := &Chat{}
	err := s.db.Get(chat, "SELECT * FROM chats WHERE id = ?", id)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return chat, nil
}

// ListChats retrieves all chats ordered by most recently updated
func (s *Store) ListChats() ([]*Chat, error) {
	chats := []*Chat{}
	err := s.db.Select(&chats, "SELECT * FROM chats ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// UpdateChatTitle updates the title of a chat
func (s *Store) UpdateChatTitle(id, title string) error {
	_, err := s.db.Exec(`
		UPDATE chats SET title = ?, updated_at = ? WHERE id = ?
	`, title, time.Now(), id)
	return err
}

// DeleteChat deletes a chat and all its messages (CASCADE)
func (s *Store) DeleteChat(id string) error {
	_, err := s.db.Exec("DELETE FROM chats WHERE id = ?", id)
	return err
}

// TouchChat updates the updated_at timestamp of a chat
func (s *Store) TouchChat(id string) error {
	_, err := s.db.Exec(`
		UPDATE chats SET updated_at = ? WHERE id = ?
	`, time.Now(), id)
	return err
}
