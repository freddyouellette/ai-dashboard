package chats_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/user_scoped_repository"
)

type ChatsRepository struct {
	*user_scoped_repository.UserScopedRepository[*models.Chat]
}

func NewChatsRepository(
	userScopedRepository *user_scoped_repository.UserScopedRepository[*models.Chat],
) *ChatsRepository {
	return &ChatsRepository{
		UserScopedRepository: userScopedRepository,
	}
}
