package repository

import (
	"context"
	"database/sql"
	"poshta/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) (int64, error)
	Delete(ctx context.Context, chatID int64) error
	GetByID(ctx context.Context, chatID int64) (*models.Chat, error)
	GetByUserID(ctx context.Context, userID int64) ([]models.Chat, error)
	GetByUsersID(ctx context.Context, user1ID, user2ID int64) (*models.Chat, error)
	GetMessages(ctx context.Context, chatID int64) ([]models.Message, error)
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
func (c *chatRepository) GetByID(ctx context.Context, chatID int64) (*models.Chat, error) {
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
func (c *chatRepository) Create(ctx context.Context, chat *models.Chat) (int64, error) {
	query := `
		INSERT INTO chats ( user1_id, user2_id, created_at )
		VALUES (?, ?, NOW())
	`
	result, err := c.db.ExecContext(ctx, query,  chat.User1ID, chat.User2ID)
	
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Delete implements ChatRepository.
func (c *chatRepository) Delete(ctx context.Context, chatID int64) error {
	panic("unimplemented")
}


// GetByUserID implements ChatRepository.
func (c *chatRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Chat, error) {
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
	var chats []models.Chat
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

func (c *chatRepository) GetByUsersID(ctx context.Context, user1ID, user2ID int64) (*models.Chat, error) {
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

func (c *chatRepository) GetMessages(ctx context.Context, chatID int64) ([]models.Message, error) {
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
    
    var messages []models.Message
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