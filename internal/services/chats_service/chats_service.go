package chats_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/services/user_scoped_service"
	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

type MessagesService interface {
	GetAllPaginated(userId uint, options *models.GetMessagesOptions) (*models.MessagesDTO, error)
	CreateWithUserId(entity *models.Message, userId uint) (*models.Message, error)
	GetByIdAndUserId(id uint, userId uint) (*models.Message, error)
	UpdateWithUserId(entity *models.Message, userId uint) (*models.Message, error)
}

type BotService interface {
	GetByIdAndUserId(id uint, userId uint) (*models.Bot, error)
}

type AiApi interface {
	GetResponse(aiModel string, randomness float64, messages []*models.Message) (*models.Message, error)
}

type ChatsService struct {
	*user_scoped_service.UserScopedService[*models.Chat]
	botService      BotService
	messagesService MessagesService
	aiApis          map[string]plugin_models.AiApiPlugin
}

func NewChatsService(
	entityService *user_scoped_service.UserScopedService[*models.Chat],
	botService BotService,
	messagesService MessagesService,
	aiApis map[string]plugin_models.AiApiPlugin,
) *ChatsService {
	return &ChatsService{
		UserScopedService: entityService,
		botService:        botService,
		messagesService:   messagesService,
		aiApis:            aiApis,
	}
}

var (
	ErrGettingChatById     = errors.New("error getting chat by id")
	ErrGettingMessageById  = errors.New("error getting message by id")
	ErrGettingBotById      = errors.New("error getting bot by id")
	ErrGettingChatMessages = errors.New("error getting chat messages")
)

func (s *ChatsService) GetChatResponse(userId uint, chatId uint) (*models.Message, error) {
	var chat *models.Chat
	var bot *models.Bot
	var err error

	chat, err = s.UserScopedService.GetByIdAndUserId(chatId, userId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingChatById, chatId, err)
	}
	bot, err = s.botService.GetByIdAndUserId(chat.BotID, userId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingBotById, chat.BotID, err)
	}
	messagesDTO, err := s.messagesService.GetAllPaginated(userId, &models.GetMessagesOptions{
		ChatID:  chatId,
		PerPage: 50,
		Page:    1,
	})
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingChatMessages, chatId, err)
	}

	// Add bot name, bot personality, and user history to list of messages to be sent
	requestMessages := make([]*plugin_models.ChatCompletionMessage, 0)

	if bot.SendName {
		requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
			Content: "Your name is " + bot.Name + ".",
			Role:    plugin_models.ChatCompletionRoleSystem,
		})
	}

	if bot.Personality != "" {
		requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
			Content: "Your personality: " + bot.Personality,
			Role:    plugin_models.ChatCompletionRoleSystem,
		})
	}

	if bot.UserHistory != "" {
		requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
			Content: "Information about me: " + bot.UserHistory,
			Role:    plugin_models.ChatCompletionRoleSystem,
		})
	}

	// NEWEST, NEWER, NEW, OLD, OLDER, OLDEST
	if len(messagesDTO.Messages) != 0 {
		// Add previous messages to list only if they are within the memory duration
		for i := len(messagesDTO.Messages) - 1; i >= 1; i-- {
			if messagesDTO.Messages[i].CreatedAt.After(time.Now().Add(-(chat.MemoryDuration * time.Second))) {
				message := messagesDTO.Messages[i]
				requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
					Content: message.Text,
					Role:    plugin_models.ChatCompletionRoleUser,
				})
			}
		}

		// ALWAYS add the last message to the list
		message := messagesDTO.Messages[0]
		requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
			Content: message.Text,
			Role:    plugin_models.ChatCompletionRoleUser,
		})
	}

	aiApi := s.aiApis[bot.AiApiPluginName]

	chatCompletionResponse, err := aiApi.CompleteChat(&plugin_models.ChatCompletionRequest{
		Model:       bot.AiModel,
		Temperature: bot.Randomness,
		Messages:    requestMessages,
	})
	if err != nil {
		return nil, err
	}

	// add this response to DB
	responseMessage := &models.Message{
		ChatID: chatId,
		Text:   chatCompletionResponse.Message.Content,
		Role:   models.MESSAGE_ROLE_BOT,
	}
	responseMessage, err = s.messagesService.CreateWithUserId(responseMessage, userId)
	if err != nil {
		return nil, err
	}

	return responseMessage, nil
}

func (s *ChatsService) GetMessageCorrection(userId uint, messageId uint) (*models.Message, error) {
	var chat *models.Chat
	var bot *models.Bot
	var err error

	message, err := s.messagesService.GetByIdAndUserId(messageId, userId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingMessageById, messageId, err)
	}
	chat, err = s.UserScopedService.GetByIdAndUserId(message.ChatID, userId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingChatById, message.ChatID, err)
	}
	bot, err = s.botService.GetByIdAndUserId(chat.BotID, userId)
	if err != nil {
		return nil, fmt.Errorf("%w (ID %d): %v", ErrGettingBotById, chat.BotID, err)
	}

	if bot.CorrectionPrompt == "" {
		return nil, errors.New("bot does not have a correction prompt")
	}

	requestMessages := make([]*plugin_models.ChatCompletionMessage, 0)

	requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
		Content: bot.CorrectionPrompt,
		Role:    plugin_models.ChatCompletionRoleSystem,
	})

	requestMessages = append(requestMessages, &plugin_models.ChatCompletionMessage{
		Content: message.Text,
		Role:    plugin_models.ChatCompletionRoleUser,
	})

	aiApi := s.aiApis[bot.AiApiPluginName]
	responseMessage, err := aiApi.CompleteChat(&plugin_models.ChatCompletionRequest{
		Model:       bot.AiModel,
		Temperature: bot.Randomness,
		Messages:    requestMessages,
	})
	if err != nil {
		return nil, err
	}

	message.Correction = responseMessage.Message.Content
	message, err = s.messagesService.UpdateWithUserId(message, userId)
	if err != nil {
		return nil, err
	}

	return message, nil
}
