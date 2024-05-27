package users_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/entity_repository"
	"gorm.io/gorm"
)

type UsersRepository struct {
	*entity_repository.EntityRepository[models.User]
	db *gorm.DB
}

func NewUsersRepository(
	entityRepository *entity_repository.EntityRepository[models.User],
	db *gorm.DB,
) *UsersRepository {
	return &UsersRepository{
		EntityRepository: entityRepository,
		db:               db,
	}
}

func (r *UsersRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
