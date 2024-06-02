package messages_service

import (
	"fmt"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/user_scoped_service"
)

type EventDispatcher interface {
	Dispatch(event models.EventType, payload interface{})
}

type MessagesRepository interface {
	GetAllPaginated(userId uint, options *models.GetMessagesOptions) (*models.MessagesDTO, error)
}

type MessagesService struct {
	*user_scoped_service.UserScopedService[*models.Message]
	messagesRepository MessagesRepository
	eventDispatcher    EventDispatcher
}

func NewMessagesService(
	entityService *user_scoped_service.UserScopedService[*models.Message],
	messagesRepository MessagesRepository,
	eventDispatcher EventDispatcher,
) *MessagesService {
	return &MessagesService{
		UserScopedService:  entityService,
		messagesRepository: messagesRepository,
		eventDispatcher:    eventDispatcher,
	}
}

func (s *MessagesService) Create(entity *models.Message) (*models.Message, error) {
	message, err := s.EntityService.Create(entity)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrRepository, err.Error())
	}

	s.eventDispatcher.Dispatch(models.EVENT_TYPE_MESSAGE_CREATED, &models.MessageCreated{
		Message: message,
	})

	return message, nil
}

func (s *MessagesService) GetAllPaginated(userId uint, options *models.GetMessagesOptions) (*models.MessagesDTO, error) {
	return s.messagesRepository.GetAllPaginated(userId, options)
}
