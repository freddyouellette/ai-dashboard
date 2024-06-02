package users_service

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type UsersRepository interface {
	GetByEmail(email string) (*models.User, error)
}

type UsersService struct {
	*entity_service.EntityService[*models.User]
	usersRepository UsersRepository
}

func NewUsersService(
	entityService *entity_service.EntityService[*models.User],
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		EntityService:   entityService,
		usersRepository: usersRepository,
	}
}

func (s *UsersService) Create(entity *models.User) (*models.User, error) {
	return s.EntityService.Create(entity)
}

func (s *UsersService) GetByEmail(email string) (*models.User, error) {
	return s.usersRepository.GetByEmail(email)
}
