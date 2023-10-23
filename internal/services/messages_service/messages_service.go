package messages_service

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type MessagesRepository interface {
	GetByChatId(chatId uint) ([]*models.Message, error)
}

type MessagesService struct {
	*entity_service.EntityService[models.Message]
	messagesRepository MessagesRepository
}

func NewMessagesService(
	entityService *entity_service.EntityService[models.Message],
	messagesRepository MessagesRepository,
) *MessagesService {
	return &MessagesService{
		EntityService:      entityService,
		messagesRepository: messagesRepository,
	}
}

func (s *MessagesService) Create(entity *models.Message) (*models.Message, error) {
	return s.EntityService.Create(entity)
}

func (s *MessagesService) GetChatMessages(chatId uint) ([]*models.Message, error) {
	entities, err := s.messagesRepository.GetByChatId(chatId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity_service.ErrRepository, err.Error())
	}
	return entities, nil
}
