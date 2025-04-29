package models

import "time"

type Message struct {
	ID        int64       `json:"id" db:"id"`
	ChatID    string       `json:"chat_id" db:"chat_id"`
	SenderID  string       `json:"sender_id" db:"sender_id"`
	SenderName string	  `json:"sender_name" db:"sender_name"`
	Content   string      `json:"content" db:"content"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	ExpiredAt time.Time   `json:"expired_at" db:"expired_at"`
	Readed 	  bool	      `json:"readed" db:"readed"`
	EncryptedKey string    `json:"encrypted_key" db:"encrypted_key"`
}

