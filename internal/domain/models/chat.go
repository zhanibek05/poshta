package models

import (
	"time"
)

type Chat struct {			
	ID        string       			`json:"id" db:"id"`
	User1ID   string     		 	`json:"user1_id" db:"user1_id"`
	User2ID   string       			`json:"user2_id" db:"user2_id"`
	CreatedAt time.Time           	`json:"created_at" db:"created_at"`

}