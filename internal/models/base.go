package models

import "gorm.io/gorm"

type BaseEntityInterface interface {
	SetID(id uint)
	GetID() uint
	IsNil() bool
}

type BaseEntity struct {
	gorm.Model
}

func (e *BaseEntity) IsNil() bool {
	return e == nil
}

func (e *BaseEntity) SetID(id uint) {
	e.ID = id
}

func (e *BaseEntity) GetID() uint {
	return e.ID
}
