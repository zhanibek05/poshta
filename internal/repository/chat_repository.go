package repository

import (
	"context"
	"database/sql"
	"poshta/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) (string, error)
	Delete(ctx context.Context, chatID string) error
	GetByID(ctx context.Context, chatID string) (*models.Chat, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Chat, error)
	GetByUsersID(ctx context.Context, user1ID, user2ID string) (*models.Chat, error)
	GetMessages(ctx context.Context, chatID string) ([]models.Message, error)
}

type chatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

// GetByID implements ChatRepository.
func (c *chatRepository) GetByID(ctx context.Context, chatID string) (*models.Chat, error) {
	query := `
		SELECT id, user1_id, user2_id, created_at
		FROM chats
		WHERE id = ?	
	`
	row := c.db.QueryRowContext(ctx, query, chatID)
	var chat models.Chat
	if err := row.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No chat found
		}
		return nil, err // Other error
	}
	return &chat, nil
}


// Create implements ChatRepository.
func (c *chatRepository) Create(ctx context.Context, chat *models.Chat) (string, error) {
	chatID := uuid.New().String()
	query := `
		INSERT INTO chats (id, user1_id, user2_id, created_at )
		VALUES (?, ?, ?, NOW())
	`
	_, err := c.db.ExecContext(ctx, query,  
		chatID,
		chat.User1ID, 
		chat.User2ID)
	
	if err != nil {
		return "", err
	}

	return chatID, nil
}

// Delete implements ChatRepository.
func (c *chatRepository) Delete(ctx context.Context, chatID string) error {
	panic("unimplemented")
}


// GetByUserID implements ChatRepository.
func (c *chatRepository) GetByUserID(ctx context.Context, userID string) ([]models.Chat, error) {
	query := `
		SELECT id, user1_id, user2_id, created_at
		FROM chats
		WHERE user1_id = ? OR user2_id = ?
	`
	rows, err := c.db.QueryContext(ctx, query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make( []models.Chat, 0)
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

// get chat by user1 and user2 id

func (c *chatRepository) GetByUsersID(ctx context.Context, user1ID, user2ID string) (*models.Chat, error) {
	query := `
		SELECT id, user1_id, user2_id, created_at
		FROM chats
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)
	`
	row := c.db.QueryRowContext(ctx, query, user1ID, user2ID, user2ID, user1ID)
	var chat models.Chat
	if err := row.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt); err != nil {
		return nil, err
	}
	return &chat, nil
}

func (c *chatRepository) GetMessages(ctx context.Context, chatID string) ([]models.Message, error) {
    query := `
        SELECT id, chat_id, sender_id, sender_name, content, created_at, readed
        FROM messages
        WHERE chat_id = ?
        ORDER BY created_at ASC
    `
    rows, err := c.db.QueryContext(ctx, query, chatID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    messages := make([]models.Message, 0)
    for rows.Next() {
        var message models.Message
        if err := rows.Scan(
            &message.ID, 
            &message.ChatID, 
            &message.SenderID, 
			&message.SenderName,
            &message.Content, 
            &message.CreatedAt,
            &message.Readed,
        ); err != nil {
            return nil, err
        }
        messages = append(messages, message)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return messages, nil
}