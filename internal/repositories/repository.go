package repositories

import "gorm.io/gorm"

type Repository[e any] struct {
	db *gorm.DB
}

func NewRepository[e any](db *gorm.DB) *Repository[e] {
	return &Repository[e]{
		db: db,
	}
}

func (r *Repository[e]) GetAll() ([]e, error) {
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

func (r *Repository[e]) GetByID(id uint) (e, error) {
	var entity e
	result := r.db.First(&entity, id)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *Repository[e]) Create(entity e) (e, error) {
	result := r.db.Create(&entity)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *Repository[e]) Update(entity e) (e, error) {
	result := r.db.Save(&entity)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *Repository[e]) Delete(entity e) error {
	result := r.db.Delete(&entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
