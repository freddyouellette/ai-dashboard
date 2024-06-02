package models

import (
	"time"
)

type Bot struct {
	BaseEntity
	UserScopedEntity
	Name             string     `json:"name"`
	UserId           uint       `json:"user_id"`
	Description      string     `json:"description"`
	SendName         bool       `json:"send_name"`
	AiApiPluginName  string     `json:"ai_api_plugin_name"`
	AiModel          string     `json:"model"`
	Personality      string     `json:"personality"`
	CorrectionPrompt string     `json:"correction_prompt"`
	UserHistory      string     `json:"user_history"`
	Randomness       float64    `json:"randomness"`
	LastMessageAt    *time.Time `json:"last_message_at"`
}
