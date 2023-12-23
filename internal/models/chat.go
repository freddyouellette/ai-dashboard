package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name           string        `json:"name"`
	BotID          uint          `json:"bot_id"`
	MemoryDuration time.Duration `json:"memory_duration"`
	LastMessageAt  *time.Time    `json:"last_message_at"`
}
