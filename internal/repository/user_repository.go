package repository

import (
	"context"
	"database/sql"
	"poshta/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (int64, error)
	Update(ctx context.Context, user *models.User) error
	GetUserPublicKey(ctx context.Context, userID int64) (string, error)
}


type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}


func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at, public_key FROM users WHERE id = ?`
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return user, nil
}


func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = ?`
	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return user, nil
}


func (r *userRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `
		INSERT INTO users (username, email, password, public_key, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.PublicKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, password = ?, updated_at = NOW() 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r * userRepository) GetUserPublicKey(ctx context.Context, userID int64) (string, error) {
	query := `SELECT public_key FROM users WHERE id = ?`
	var publicKey string
	err := r.db.GetContext(ctx, &publicKey, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // User not found
		}
		return "", err
	}
	return publicKey, nil
}