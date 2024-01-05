package ai_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AiApi struct {
	client      HttpClient
	maxTokens   int
	accessToken string
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
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
	accessToken string,
) *AiApi {
	return &AiApi{
		client:      client,
		maxTokens:   maxTokens,
		accessToken: accessToken,
	}
}

var (
	ErrBadResponse = errors.New("received bad response")
)

func (api *AiApi) GetResponse(aiModel string, randomness float64, messages []*models.Message) (*models.Message, error) {
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

	requestBody := chatRequest{
		Model:       aiModel,
		Messages:    requestMessages,
		MaxTokens:   api.maxTokens,
		Temperature: randomness,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestReader := bytes.NewBuffer(jsonBody)

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", requestReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+api.accessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody chatResponse
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

func (api *AiApi) GetBotModels() ([]*models.BotModel, error) {
	request, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+api.accessToken)
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
	for _, model := range responseBody.Data {
		botModels = append(botModels, &models.BotModel{
			ID:        model.Id,
			CreatedAt: time.Unix(model.Created, 0),
		})
	}

	return botModels, nil
}
