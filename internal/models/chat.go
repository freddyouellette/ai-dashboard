package models

import (
	"time"
)

type Chat struct {
	BaseEntity
	UserScopedEntity
	Name           string        `json:"name"`
	BotID          uint          `json:"bot_id"`
	MemoryDuration time.Duration `json:"memory_duration"`
	LastMessageAt  *time.Time    `json:"last_message_at"`
}
