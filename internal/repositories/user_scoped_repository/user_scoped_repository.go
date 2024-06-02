package user_scoped_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"gorm.io/gorm"
)

type EntityRepository[e models.UserScopedEntityInterface] interface {
	GetAll() ([]e, error)
	GetById(id uint) (e, error)
	Create(entity e) (e, error)
	Update(entity e) (e, error)
	Delete(entity e) error
}

type UserScopedRepository[e models.UserScopedEntityInterface] struct {
	EntityRepository[e]
	blank e
	db    *gorm.DB
}

func NewUserScopedRepository[e models.UserScopedEntityInterface](
	entityRepository EntityRepository[e],
	db *gorm.DB,
) *UserScopedRepository[e] {
	var blank e
	return &UserScopedRepository[e]{
		blank:            blank,
		EntityRepository: entityRepository,
		db:               db,
	}
}

func (r *UserScopedRepository[e]) GetAllWithUserId(userId uint) ([]e, error) {
	var entities []e
	result := r.db.Where("user_id = ?", userId).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(entities) == 0 {
		// if there are no entities, return empty slice instead of nil...
		return make([]e, 0), nil
	}
	return entities, nil
}

func (r *UserScopedRepository[e]) GetByIdAndUserId(id uint, userId uint) (e, error) {
	var entity e
	result := r.db.Where("user_id = ?", userId).First(&entity, id)
	if result.Error != nil {
		return r.blank, result.Error
	}
	return entity, nil
}

func (r *UserScopedRepository[e]) CreateWithUserId(entity e, userId uint) (e, error) {
	(entity).SetUserId(userId)
	result := r.db.Create(entity)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *UserScopedRepository[e]) UpdateWithUserId(entity e, userId uint) (e, error) {
	old, err := r.GetByIdAndUserId((entity).GetID(), userId)
	if err != nil {
		return r.blank, err
	}
	if old.IsNil() {
		return r.blank, models.ErrResourceNotFound
	}
	result := r.db.Save(entity)
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}

func (r *UserScopedRepository[e]) Delete(entity e) error {
	result := r.db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
