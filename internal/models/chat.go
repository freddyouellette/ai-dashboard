package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name  string `json:"name"`
	BotID uint   `json:"bot_id"`
}
