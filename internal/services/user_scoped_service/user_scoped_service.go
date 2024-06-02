package user_scoped_service

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type UserScopedRepository[e models.UserScopedEntityInterface] interface {
	Delete(entity e) error
	GetAllWithUserId(userId uint) ([]e, error)
	CreateWithUserId(entity e, userId uint) (e, error)
	UpdateWithUserId(entity e, userId uint) (e, error)
	GetByIdAndUserId(id uint, userId uint) (e, error)
}

type UserScopedService[e models.UserScopedEntityInterface] struct {
	*entity_service.EntityService[e]
	entityRepository UserScopedRepository[e]
}

func NewUserScopedService[e models.UserScopedEntityInterface](
	entityService *entity_service.EntityService[e],
	entityRepository UserScopedRepository[e],
) *UserScopedService[e] {
	return &UserScopedService[e]{
		EntityService:    entityService,
		entityRepository: entityRepository,
	}
}

func (s *UserScopedService[e]) GetAllWithUserId(userId uint) ([]e, error) {
	return s.entityRepository.GetAllWithUserId(userId)
}

func (s *UserScopedService[e]) CreateWithUserId(entity e, userId uint) (e, error) {
	return s.entityRepository.CreateWithUserId(entity, userId)
}

func (s *UserScopedService[e]) UpdateWithUserId(entity e, userId uint) (e, error) {
	return s.entityRepository.UpdateWithUserId(entity, userId)
}

func (s *UserScopedService[e]) GetByIdAndUserId(id uint, userId uint) (e, error) {
	return s.entityRepository.GetByIdAndUserId(id, userId)
}

func (s *UserScopedService[e]) DeleteFromUserId(id uint, userId uint) error {
	entity, err := s.entityRepository.GetByIdAndUserId(id, userId)
	if err != nil {
		return fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	err = s.entityRepository.Delete(entity)
	if err != nil {
		return fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	return nil
}
