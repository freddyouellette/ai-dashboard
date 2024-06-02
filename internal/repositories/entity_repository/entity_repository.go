package entity_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"gorm.io/gorm"
)

type EntityRepository[e models.BaseEntityInterface] struct {
	blank e
	db    *gorm.DB
}

func NewRepository[e models.BaseEntityInterface](db *gorm.DB) *EntityRepository[e] {
	var blank e
	return &EntityRepository[e]{
		blank: blank,
		db:    db,
	}
}

func (r *EntityRepository[e]) GetAll() ([]e, error) {
	var entities []e
	result := r.db.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(entities) == 0 {
		// if there are no entities, return empty slice instead of nil...
		return make([]e, 0), nil
	}
	return entities, nil
}

func (r *EntityRepository[e]) GetById(id uint) (e, error) {
	var entity e
	result := r.db.First(&entity, id)
	if result.Error != nil {
		return r.blank, result.Error
	}
	return entity, nil
}

func (r *EntityRepository[e]) Create(entity e) (e, error) {
	result := r.db.Create(entity)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *EntityRepository[e]) Update(entity e) (e, error) {
	old, err := r.GetById(entity.GetID())
	if err != nil {
		return r.blank, err
	}
	if old.GetID() == 0 {
		return r.blank, models.ErrResourceNotFound
	}
	result := r.db.Save(entity)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *EntityRepository[e]) Delete(entity e) error {
	result := r.db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
