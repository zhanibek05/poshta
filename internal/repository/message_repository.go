package repository

import (
	"context"
	"poshta/internal/domain/models"
	
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) (int64, error)
	GetByID(ctx context.Context, messageID int64) (*models.Message, error)
}

type messageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (m *messageRepository) Create(ctx context.Context, message *models.Message) (int64, error) {
	query := `
		INSERT INTO messages (chat_id, sender_id, sender_name, content, created_at)
		VALUES (?, ?, ?, ?, NOW())
	`
	// get sender name from user_id from user repository
	
	result, err := m.db.ExecContext(ctx, query, message.ChatID, message.SenderID, message.SenderName, message.Content)
	if err != nil {
		return 0, err
	}
	messageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return messageID, nil
}


func (m *messageRepository) GetByID(ctx context.Context, messageID int64) (*models.Message, error) {
	query := `
		SELECT id, chat_id, sender_id, sender_name, content, created_at
		FROM messages
		WHERE id = ?	
	`
	row := m.db.QueryRowContext(ctx, query, messageID)
	var message models.Message
	if err := row.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.SenderName, &message.Content, &message.CreatedAt); err != nil {
		return nil, err // Other error
	}
	return &message, nil
}