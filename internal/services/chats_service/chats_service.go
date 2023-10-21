package chats_service

import (
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type MessagesService interface {
	GetChatMessages(chatId uint) ([]*models.Message, error)
	Create(message *models.Message) (*models.Message, error)
}

type BotService interface {
	GetById(botId uint) (*models.Bot, error)
}

type AiApi interface {
	GetResponse(aiModel string, messages []*models.Message) (*models.Message, error)
}

type ChatsService struct {
	*entity_service.EntityService[models.Chat]
	botService      BotService
	messagesService MessagesService
	aiApi           AiApi
}

func NewChatsService(
	entityService *entity_service.EntityService[models.Chat],
	botService BotService,
	messagesService MessagesService,
	aiApi AiApi,
) *ChatsService {
	return &ChatsService{
		EntityService:   entityService,
		botService:      botService,
		messagesService: messagesService,
		aiApi:           aiApi,
	}
}

func (s *ChatsService) GetChatResponse(chatId uint) (*models.Message, error) {
	chat, err := s.EntityService.GetById(chatId)
	if err != nil {
		return nil, err
	}
	bot, err := s.botService.GetById(chat.BotID)
	if err != nil {
		return nil, err
	}
	messages, err := s.messagesService.GetChatMessages(chatId)
	if err != nil {
		return nil, err
	}

	// Add bot name, bot personality, and user history to list of messages to be sent
	requestMessages := make([]*models.Message, 0)

	requestMessages = append(requestMessages, &models.Message{
		Text: "Your name is " + bot.Name + ".",
		Role: models.MESSAGE_ROLE_SYSTEM,
	})

	if bot.Description != "" {
		requestMessages = append(requestMessages, &models.Message{
			Text: "Your personality: " + bot.Personality,
			Role: models.MESSAGE_ROLE_SYSTEM,
		})
	}

	if bot.UserHistory != "" {
		requestMessages = append(requestMessages, &models.Message{
			Text: "Information about me: " + bot.UserHistory,
			Role: models.MESSAGE_ROLE_SYSTEM,
		})
	}

	// Add previous messages to list
	// TODO: Limit this
	requestMessages = append(requestMessages, messages...)

	responseMessage, err := s.aiApi.GetResponse(bot.AiModel, requestMessages)
	if err != nil {
		return nil, err
	}

	// add this response to DB
	responseMessage.ChatID = chatId
	responseMessage.Role = models.MESSAGE_ROLE_BOT
	responseMessage, err = s.messagesService.Create(responseMessage)
	if err != nil {
		return nil, err
	}

	return responseMessage, nil
}
