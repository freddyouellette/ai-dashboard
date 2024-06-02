package bots_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/user_scoped_repository"
)

type BotsRepository struct {
	*user_scoped_repository.UserScopedRepository[*models.Bot]
}

func NewBotsRepository(
	userScopedRepository *user_scoped_repository.UserScopedRepository[*models.Bot],
) *BotsRepository {
	return &BotsRepository{
		UserScopedRepository: userScopedRepository,
	}
}
