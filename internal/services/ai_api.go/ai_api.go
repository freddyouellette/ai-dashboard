package ai_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AiApi struct {
	client      HttpClient
	maxTokens   int
	randomness  float32
	chatUrl     string
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
	Temperature float32       `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewAiApi(
	client HttpClient,
	maxTokens int,
	randomness float32,
	chatUrl string,
	accessToken string,
) *AiApi {
	return &AiApi{
		client:      client,
		maxTokens:   maxTokens,
		randomness:  randomness,
		chatUrl:     chatUrl,
		accessToken: accessToken,
	}
}

var (
	ErrBadResponse = errors.New("received bad response")
)

func (api *AiApi) GetResponse(aiModel string, messages []*models.Message) (*models.Message, error) {
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
		Temperature: api.randomness,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestReader := bytes.NewBuffer(jsonBody)

	request, err := http.NewRequest(http.MethodPost, api.chatUrl, requestReader)
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
