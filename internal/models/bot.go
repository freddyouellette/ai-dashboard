package models

import "gorm.io/gorm"

type Bot struct {
	gorm.Model
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AiModel     string `json:"model"`
	Personality string `json:"personality"`
	UserHistory string `json:"user_history"`
}
