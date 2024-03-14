package ai_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

const (
	ANTHROPIC_CLAUDE_3_OPUS   = "claude-3-opus-20240229"
	ANTHROPIC_CLAUDE_3_SONNET = "claude-3-sonnet-20240229"
	ANTHROPIC_CLAUDE_3_HAIKU  = "claude-3-haiku-20240307"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AiApi struct {
	client               HttpClient
	maxTokens            int
	chatGptAccessToken   string
	anthropicAccessToken string
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatGptRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type chatGptResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type anthropicRequest struct {
	Model       string        `json:"model"`
	System      string        `json:"system"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
}

type model struct {
	Id      string `json:"id"`
	Created int64  `json:"created"`
}

type modelsResponse struct {
	Data []model `json:"data"`
}

func NewAiApi(
	client HttpClient,
	maxTokens int,
	chatGptAccessToken string,
	anthropicAccessToken string,
) *AiApi {
	return &AiApi{
		client:               client,
		maxTokens:            maxTokens,
		chatGptAccessToken:   chatGptAccessToken,
		anthropicAccessToken: anthropicAccessToken,
	}
}

var (
	ErrBadResponse = errors.New("received bad response")
)

func (api *AiApi) GetResponse(aiModel string, randomness float64, messages []*models.Message) (*models.Message, error) {
	if isAnthropic(aiModel) {
		return api.getAnthropicResponse(aiModel, randomness, messages)
	} else {
		return api.getChatGptResponse(aiModel, randomness, messages)
	}
}

func (api *AiApi) getChatGptResponse(aiModel string, randomness float64, messages []*models.Message) (*models.Message, error) {
	requestBody := chatGptRequest{
		Model:       aiModel,
		MaxTokens:   api.maxTokens,
		Temperature: randomness,
	}

	requestMessages := make([]chatMessage, 0)
	for _, message := range messages {
		var messageRole string
		switch message.Role {
		case models.MESSAGE_ROLE_USER:
			messageRole = "user"
		case models.MESSAGE_ROLE_BOT:
			messageRole = "assistant"
		case models.MESSAGE_ROLE_SYSTEM:
			messageRole = "system"
		default:
			continue
		}

		requestMessages = append(requestMessages, chatMessage{
			Role:    messageRole,
			Content: message.Text,
		})
	}

	requestBody.Messages = requestMessages

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestReader := bytes.NewBuffer(jsonBody)

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", requestReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+api.chatGptAccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody chatGptResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %s", ErrBadResponse, response.StatusCode, responseBody)
	}

	return &models.Message{
		Text: responseBody.Choices[0].Message.Content,
		Role: models.MESSAGE_ROLE_BOT,
	}, nil
}

var anthropicAttempts int = 0

func (api *AiApi) getAnthropicResponse(aiModel string, randomness float64, messages []*models.Message) (*models.Message, error) {
	requestBody := anthropicRequest{
		Model:       aiModel,
		MaxTokens:   api.maxTokens,
		Temperature: randomness,
	}

	requestMessages := make([]chatMessage, 0)
	for _, message := range messages {
		var messageRole string
		switch message.Role {
		case models.MESSAGE_ROLE_USER:
			messageRole = "user"
		case models.MESSAGE_ROLE_BOT:
			messageRole = "assistant"
		case models.MESSAGE_ROLE_SYSTEM:
			if requestBody.System == "" {
				requestBody.System = message.Text
			} else {
				requestBody.System += "\n\n" + message.Text
			}
			continue
		default:
			continue
		}

		requestMessages = append(requestMessages, chatMessage{
			Role:    messageRole,
			Content: message.Text,
		})
	}

	requestBody.Messages = requestMessages

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestReader := bytes.NewBuffer(jsonBody)

	request, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", requestReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("anthropic-version", "2023-06-01")
	request.Header.Set("x-api-key", api.anthropicAccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody anthropicResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if (response.StatusCode == 529 || response.StatusCode == 429) && anthropicAttempts < 3 {
		anthropicAttempts++
		time.Sleep(time.Duration(math.Pow(3.0, float64(anthropicAttempts))) * time.Second)
		return api.getAnthropicResponse(aiModel, randomness, messages)
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %s", ErrBadResponse, response.StatusCode, responseBody)
	}

	return &models.Message{
		Text: responseBody.Content[0].Text,
		Role: models.MESSAGE_ROLE_BOT,
	}, nil
}

func isAnthropic(model string) bool {
	return model == ANTHROPIC_CLAUDE_3_OPUS || model == ANTHROPIC_CLAUDE_3_SONNET
}

func (api *AiApi) GetBotModels() ([]*models.BotModel, error) {
	request, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+api.chatGptAccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody modelsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %v", ErrBadResponse, response.StatusCode, responseBody)
	}

	botModels := make([]*models.BotModel, 0)

	botModels = append(botModels, &models.BotModel{
		ID:        ANTHROPIC_CLAUDE_3_OPUS,
		CreatedAt: time.Date(2024, 02, 29, 0, 0, 0, 3, time.UTC),
	})
	botModels = append(botModels, &models.BotModel{
		ID:        ANTHROPIC_CLAUDE_3_SONNET,
		CreatedAt: time.Date(2024, 02, 29, 0, 0, 0, 2, time.UTC),
	})
	botModels = append(botModels, &models.BotModel{
		ID:        ANTHROPIC_CLAUDE_3_HAIKU,
		CreatedAt: time.Date(2024, 02, 29, 0, 0, 0, 1, time.UTC),
	})

	for _, model := range responseBody.Data {
		botModels = append(botModels, &models.BotModel{
			ID:        model.Id,
			CreatedAt: time.Unix(model.Created, 0),
		})
	}

	return botModels, nil
}
