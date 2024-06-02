package models

import (
	"time"
)

type Chat struct {
	BaseEntity
	UserScopedEntity
	Name           string        `json:"name"`
	UserId         uint          `json:"user_id"`
	BotID          uint          `json:"bot_id"`
	MemoryDuration time.Duration `json:"memory_duration"`
	LastMessageAt  *time.Time    `json:"last_message_at"`
}
