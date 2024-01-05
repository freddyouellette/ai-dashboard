package models

import "time"

type BotModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
