package messages_repository

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/entity_repository"
	"gorm.io/gorm"
)

type MessagesRepository struct {
	*entity_repository.EntityRepository[models.Message]
	db *gorm.DB
}

func NewMessagesRepository(
	entityRepository *entity_repository.EntityRepository[models.Message],
	db *gorm.DB,
) *MessagesRepository {
	return &MessagesRepository{
		EntityRepository: entityRepository,
		db:               db,
	}
}

func (r *MessagesRepository) GetByChatId(chatId uint) ([]*models.Message, error) {
	var messages []*models.Message
	result := r.db.Where("chat_id = ?", chatId).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(messages) == 0 {
		// if there are no entities, return empty slice instead of nil...
		return make([]*models.Message, 0), nil
	}
	return messages, nil
}
