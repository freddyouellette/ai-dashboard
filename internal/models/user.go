package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID    string `json:"id"`
	Email string `json:"email"`
}
