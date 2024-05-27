package users_service

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type UsersRepository interface {
	GetByEmail(email string) (*models.User, error)
}

type UsersService struct {
	*entity_service.EntityService[models.User]
	usersRepository UsersRepository
}

func NewUsersService(
	entityService *entity_service.EntityService[models.User],
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		EntityService:   entityService,
		usersRepository: usersRepository,
	}
}

func (s *UsersService) Create(entity *models.User) (*models.User, error) {
	user, err := s.EntityService.Create(entity)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity_service.ErrRepository, err.Error())
	}

	return user, nil
}

func (s *UsersService) GetByEmail(email string) (*models.User, error) {
	user, err := s.usersRepository.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity_service.ErrRepository, err.Error())
	}

	return user, nil
}
