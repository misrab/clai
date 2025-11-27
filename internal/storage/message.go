package storage

import (
	"time"
)

// Message represents a single message in a chat
type Message struct {
	ID        string    `json:"id" db:"id"`
	ChatID    string    `json:"chat_id" db:"chat_id"`
	Role      string    `json:"role" db:"role"` // "user" or "assistant"
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateMessage creates a new message in the database
func (s *Store) CreateMessage(msg *Message) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert message
	_, err = tx.Exec(`
		INSERT INTO messages (id, chat_id, role, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, msg.ID, msg.ChatID, msg.Role, msg.Content, msg.CreatedAt)
	if err != nil {
		return err
	}

	// Update chat's updated_at timestamp
	_, err = tx.Exec(`
		UPDATE chats SET updated_at = ? WHERE id = ?
	`, time.Now(), msg.ChatID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetMessages retrieves all messages for a chat, ordered by creation time
func (s *Store) GetMessages(chatID string) ([]*Message, error) {
	messages := []*Message{}
	err := s.db.Select(&messages, `
		SELECT * FROM messages
		WHERE chat_id = ?
		ORDER BY created_at ASC
	`, chatID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// DeleteMessage deletes a single message
func (s *Store) DeleteMessage(id string) error {
	_, err := s.db.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}

// GetChatWithMessages retrieves a chat with all its messages
func (s *Store) GetChatWithMessages(chatID string) (*Chat, []*Message, error) {
	chat, err := s.GetChat(chatID)
	if err != nil {
		return nil, nil, err
	}
	if chat == nil {
		return nil, nil, nil
	}

	messages, err := s.GetMessages(chatID)
	if err != nil {
		return nil, nil, err
	}

	return chat, messages, nil
}
