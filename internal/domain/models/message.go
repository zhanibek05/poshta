package models

import "time"

type Message struct {
	ID        int       `json:"id" db:"id"`
	ChatID    int       `json:"chat_id" db:"chat_id"`
	SenderID  int       `json:"sender_id" db:"sender_id"`
	Content   string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
	Readed 	  bool	    `json:"readed" db:"readed"`
}

