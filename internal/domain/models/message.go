package models

import "time"

type Message struct {
	ID        int64       `json:"id" db:"id"`
	ChatID    int64       `json:"chat_id" db:"chat_id"`
	SenderID  int64       `json:"sender_id" db:"sender_id"`
	SenderName string	  `json:"sender_name" db:"sender_name"`
	Content   string      `json:"content" db:"content"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	ExpiredAt time.Time   `json:"expired_at" db:"expired_at"`
	Readed 	  bool	      `json:"readed" db:"readed"`
	EncryptedKey string    `json:"encrypted_key" db:"encrypted_key"`
	EncryptedKeySender string `json:"encrypted_key_sender" db:"encrypted_key_sender"`
}

