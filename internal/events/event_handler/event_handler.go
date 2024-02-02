package event_handler

import "github.com/freddyouellette/ai-dashboard/internal/models"

type ChatsService interface {
	GetById(id uint) (*models.Chat, error)
	Update(entity *models.Chat) (*models.Chat, error)
}

type BotsService interface {
	GetById(id uint) (*models.Bot, error)
	Update(entity *models.Bot) (*models.Bot, error)
}

type EventHandler struct {
	chatsService ChatsService
	botsService  BotsService
}

func NewEventHandler(chatsService ChatsService, botsService BotsService) *EventHandler {
	return &EventHandler{
		chatsService: chatsService,
		botsService:  botsService,
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

	bot, err := h.botsService.GetById(chat.BotID)
	if err != nil {
		return
	}

	bot.LastMessageAt = &event.Message.CreatedAt
	_, err = h.botsService.Update(bot)
	if err != nil {
		return
	}
}
