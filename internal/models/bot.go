package models

import (
	"time"

	"gorm.io/gorm"
)

type Bot struct {
	gorm.Model
	Name             string     `json:"name"`
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
