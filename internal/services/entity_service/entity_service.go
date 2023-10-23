package entity_service

import (
	"fmt"
)

type EntityRepository[e any] interface {
	GetAll() ([]*e, error)
	GetByID(id uint) (*e, error)
	Create(entity *e) (*e, error)
	Update(entity *e) (*e, error)
	Delete(entity *e) error
}

type EntityService[e any] struct {
	entityRepository EntityRepository[e]
}

var (
	ErrRepository = fmt.Errorf("repository error")
)

func NewEntityService[e any](entityRepository EntityRepository[e]) *EntityService[e] {
	return &EntityService[e]{
		entityRepository: entityRepository,
	}
}

func (s *EntityService[e]) GetAll() ([]*e, error) {
	entities, err := s.entityRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return entities, nil
}

func (s *EntityService[e]) Create(entity *e) (*e, error) {
	entity, err := s.entityRepository.Create(entity)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return entity, nil
}

func (s *EntityService[e]) Update(entity *e) (*e, error) {
	entity, err := s.entityRepository.Update(entity)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return entity, nil
}

func (s *EntityService[e]) GetById(id uint) (*e, error) {
	entity, err := s.entityRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return entity, nil
}

func (s *EntityService[e]) Delete(id uint) error {
	entity, err := s.entityRepository.GetByID(id)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	err = s.entityRepository.Delete(entity)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}
	return nil
}
