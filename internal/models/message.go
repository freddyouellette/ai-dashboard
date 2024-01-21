package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ChatID     uint        `json:"chat_id"`
	Text       string      `json:"text"`
	Correction string      `json:"correction"`
	Role       MessageRole `json:"role"`
}
