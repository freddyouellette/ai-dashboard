package messages_service

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type EventDispatcher interface {
	Dispatch(event models.EventType, payload interface{})
}

type MessagesRepository interface {
	GetAllPaginated(options *models.GetMessagesOptions) (*models.MessagesDTO, error)
}

type MessagesService struct {
	*entity_service.EntityService[models.Message]
	messagesRepository MessagesRepository
	eventDispatcher    EventDispatcher
}

func NewMessagesService(
	entityService *entity_service.EntityService[models.Message],
	messagesRepository MessagesRepository,
	eventDispatcher EventDispatcher,
) *MessagesService {
	return &MessagesService{
		EntityService:      entityService,
		messagesRepository: messagesRepository,
		eventDispatcher:    eventDispatcher,
	}
}

func (s *MessagesService) Create(entity *models.Message) (*models.Message, error) {
	message, err := s.EntityService.Create(entity)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity_service.ErrRepository, err.Error())
	}

	s.eventDispatcher.Dispatch(models.EVENT_TYPE_MESSAGE_CREATED, &models.MessageCreated{
		Message: message,
	})

	return message, nil
}

func (s *MessagesService) GetAllPaginated(options *models.GetMessagesOptions) (*models.MessagesDTO, error) {
	messagesDTO, err := s.messagesRepository.GetAllPaginated(options)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity_service.ErrRepository, err.Error())
	}
	return messagesDTO, nil
}
