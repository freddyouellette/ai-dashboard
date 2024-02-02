package chats_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
)

type MessagesService interface {
	GetAllPaginated(options *models.GetMessagesOptions) (*models.MessagesDTO, error)
	Create(message *models.Message) (*models.Message, error)
	GetById(messageId uint) (*models.Message, error)
	Update(message *models.Message) (*models.Message, error)
}

type BotService interface {
	GetById(botId uint) (*models.Bot, error)
}

type AiApi interface {
	GetResponse(aiModel string, randomness float64, messages []*models.Message) (*models.Message, error)
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

var (
	ErrGettingChatById     = errors.New("error getting chat by id")
	ErrGettingMessageById  = errors.New("error getting message by id")
	ErrGettingBotById      = errors.New("error getting bot by id")
	ErrGettingChatMessages = errors.New("error getting chat messages")
)

func (s *ChatsService) GetChatResponse(chatId uint) (*models.Message, error) {
	var chat *models.Chat
	var bot *models.Bot
	var err error

	chat, err = s.EntityService.GetById(chatId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingChatById, chatId, err)
	}
	bot, err = s.botService.GetById(chat.BotID)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingBotById, chat.BotID, err)
	}
	messagesDTO, err := s.messagesService.GetAllPaginated(&models.GetMessagesOptions{
		ChatID:  chatId,
		PerPage: 50,
		Page:    1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingChatMessages, chatId, err)
	}

	// Add bot name, bot personality, and user history to list of messages to be sent
	requestMessages := make([]*models.Message, 0)

	if bot.SendName {
		requestMessages = append(requestMessages, &models.Message{
			Text: "Your name is " + bot.Name + ".",
			Role: models.MESSAGE_ROLE_SYSTEM,
		})
	}

	if bot.Personality != "" {
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

	// NEWEST, NEWER, NEW, OLD, OLDER, OLDEST
	if len(messagesDTO.Messages) != 0 {
		// Add previous messages to list only if they are within the memory duration
		for i := len(messagesDTO.Messages) - 1; i >= 1; i-- {
			if messagesDTO.Messages[i].CreatedAt.After(time.Now().Add(-(chat.MemoryDuration * time.Second))) {
				requestMessages = append(requestMessages, messagesDTO.Messages[i])
			}
		}

		// ALWAYS add the last message to the list
		requestMessages = append(requestMessages, messagesDTO.Messages[0])
	}

	var responseMessage *models.Message
	responseMessage, err = s.aiApi.GetResponse(bot.AiModel, bot.Randomness, requestMessages)
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

func (s *ChatsService) GetMessageCorrection(messageId uint) (*models.Message, error) {
	var chat *models.Chat
	var bot *models.Bot
	var err error

	message, err := s.messagesService.GetById(messageId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingMessageById, messageId, err)
	}
	chat, err = s.EntityService.GetById(message.ChatID)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingChatById, message.ChatID, err)
	}
	bot, err = s.botService.GetById(chat.BotID)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingBotById, chat.BotID, err)
	}

	if bot.CorrectionPrompt == "" {
		return nil, errors.New("bot does not have a correction prompt")
	}

	requestMessages := make([]*models.Message, 0)

	requestMessages = append(requestMessages, &models.Message{
		Text: bot.CorrectionPrompt,
		Role: models.MESSAGE_ROLE_SYSTEM,
	})

	requestMessages = append(requestMessages, message)

	var responseMessage *models.Message
	responseMessage, err = s.aiApi.GetResponse(bot.AiModel, bot.Randomness, requestMessages)
	if err != nil {
		return nil, err
	}

	message.Correction = responseMessage.Text
	message, err = s.messagesService.Update(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
