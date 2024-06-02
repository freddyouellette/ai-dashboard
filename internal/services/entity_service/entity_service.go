package entity_service

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type EntityRepository[e models.BaseEntityInterface] interface {
	GetAll() ([]e, error)
	GetById(id uint) (e, error)
	Create(entity e) (e, error)
	Update(entity e) (e, error)
	Delete(entity e) error
}

type EntityService[e models.BaseEntityInterface] struct {
	blank            e
	entityRepository EntityRepository[e]
}

func NewEntityService[e models.BaseEntityInterface](entityRepository EntityRepository[e]) *EntityService[e] {
	var blank e
	return &EntityService[e]{
		blank:            blank,
		entityRepository: entityRepository,
	}
}

func (s *EntityService[e]) GetAll() ([]e, error) {
	entities, err := s.entityRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	return entities, nil
}

func (s *EntityService[e]) Create(entity e) (e, error) {
	entity, err := s.entityRepository.Create(entity)
	if err != nil {
		return s.blank, fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	return entity, nil
}

func (s *EntityService[e]) Update(entity e) (e, error) {
	entity, err := s.entityRepository.Update(entity)
	if err != nil {
		return s.blank, fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	return entity, nil
}

func (s *EntityService[e]) GetById(id uint) (e, error) {
	entity, err := s.entityRepository.GetById(id)
	if err != nil {
		return s.blank, fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	return entity, nil
}

func (s *EntityService[e]) Delete(id uint) error {
	entity, err := s.entityRepository.GetById(id)
	if err != nil {
		return fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	err = s.entityRepository.Delete(entity)
	if err != nil {
		return fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}
	return nil
}
