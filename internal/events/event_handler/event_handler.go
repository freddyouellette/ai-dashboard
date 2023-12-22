package event_handler

import "github.com/freddyouellette/ai-dashboard/internal/models"

type ChatsService interface {
	GetById(id uint) (*models.Chat, error)
	Update(entity *models.Chat) (*models.Chat, error)
}

type EventHandler struct {
	chatsService ChatsService
}

func NewEventHandler(chatsService ChatsService) *EventHandler {
	return &EventHandler{
		chatsService: chatsService,
	}
}

func (h *EventHandler) HandleMessageCreatedEvent(payload interface{}) {
	event := payload.(*models.MessageCreated)

	chat, err := h.chatsService.GetById(event.Message.ChatID)
	if err != nil {
		return
	}

	chat.LastMessageAt = &event.Message.CreatedAt
	_, err = h.chatsService.Update(chat)
	if err != nil {
		return
	}
}
