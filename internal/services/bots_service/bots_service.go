package bots_service

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/user_scoped_service"
)

type BotsService struct {
	*user_scoped_service.UserScopedService[*models.Bot]
}

func NewBotsService(
	userScopedService *user_scoped_service.UserScopedService[*models.Bot],
) *BotsService {
	return &BotsService{
		UserScopedService: userScopedService,
	}
}
