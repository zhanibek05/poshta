package models

import (
	"time"
)

type Chat struct {			
	ID        int64       			`json:"id" db:"id"`
	User1ID   int64      		 	`json:"user1_id" db:"user1_id"`
	User2ID   int64       			`json:"user2_id" db:"user2_id"`
	CreatedAt time.Time           	`json:"created_at" db:"created_at"`

}