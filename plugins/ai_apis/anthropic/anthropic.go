package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

const (
	PLUGIN_ID   = "anthropic"
	PLUGIN_NAME = "Anthropic"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Anthropic struct {
	client               HttpClient
	maxTokens            int
	anthropicAccessToken string
}

type anthropicChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicChatCompletionRequest struct {
	Model       string                           `json:"model"`
	System      string                           `json:"system"`
	Messages    []anthropicChatCompletionMessage `json:"messages"`
	MaxTokens   int                              `json:"max_tokens"`
	Temperature float64                          `json:"temperature"`
}

type anthropicChatCompletionResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
}

type anthropicModel struct {
	Id      string `json:"id"`
	Created int64  `json:"created"`
}

type anthropicModelsResponse struct {
	Data []anthropicModel `json:"data"`
}

func NewAnthropic(
	client HttpClient,
	maxTokens int,
	anthropicAccessToken string,
) *Anthropic {
	return &Anthropic{
		client:               client,
		maxTokens:            maxTokens,
		anthropicAccessToken: anthropicAccessToken,
	}
}

var (
	ErrBadResponse = errors.New("received bad response")
)

func (api *Anthropic) Initialize(options *plugin_models.AiApiPluginOptions) error {
	api.client = options.Client
	maxTokens, err := strconv.Atoi(os.Getenv("ANTHROPIC_MAX_TOKENS"))
	if err != nil {
		return err
	}
	api.maxTokens = maxTokens
	api.anthropicAccessToken = os.Getenv("ANTHROPIC_ACCESS_TOKEN")
	return nil
}

var anthropicAttempts int = 0

func (api *Anthropic) CompleteChat(chatCompletionRequest *plugin_models.ChatCompletionRequest) (*plugin_models.ChatCompletionResponse, error) {
	requestBody := anthropicChatCompletionRequest{
		Model:       chatCompletionRequest.Model,
		MaxTokens:   api.maxTokens,
		Temperature: chatCompletionRequest.Temperature,
	}

	requestMessages := make([]anthropicChatCompletionMessage, 0)
	for _, message := range chatCompletionRequest.Messages {
		var messageRole string
		switch message.Role {
		case plugin_models.ChatCompletionRoleUser:
			messageRole = "user"
		case plugin_models.ChatCompletionRoleAssistant:
			messageRole = "assistant"
		case plugin_models.ChatCompletionRoleSystem:
			if requestBody.System == "" {
				requestBody.System = message.Content
			} else {
				requestBody.System += "\n\n" + message.Content
			}
			continue
		default:
			continue
		}

		requestMessages = append(requestMessages, anthropicChatCompletionMessage{
			Role:    messageRole,
			Content: message.Content,
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

	var responseBody anthropicChatCompletionResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if (response.StatusCode == 529 || response.StatusCode == 429) && anthropicAttempts < 3 {
		anthropicAttempts++
		time.Sleep(time.Duration(math.Pow(3.0, float64(anthropicAttempts))) * time.Second)
		return api.CompleteChat(chatCompletionRequest)
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %s", ErrBadResponse, response.StatusCode, responseBody)
	}

	return &plugin_models.ChatCompletionResponse{
		Message: &plugin_models.ChatCompletionMessage{
			Content: responseBody.Content[0].Text,
			Role:    plugin_models.ChatCompletionRoleAssistant,
		},
	}, nil
}

func (api *Anthropic) GetModels() (*plugin_models.GetModelsResponse, error) {
	return &plugin_models.GetModelsResponse{
		Models: []*plugin_models.AiModel{
			{
				Id:         "claude-3-opus-20240229",
				AuthorId:   PLUGIN_ID,
				AuthorName: PLUGIN_NAME,
			},
			{
				Id:         "claude-3-sonnet-20240229",
				AuthorId:   PLUGIN_ID,
				AuthorName: PLUGIN_NAME,
			},
			{
				Id:         "claude-3-haiku-20240307",
				AuthorId:   PLUGIN_ID,
				AuthorName: PLUGIN_NAME,
			},
		},
	}, nil
}

func (api *Anthropic) GetPluginId() string {
	return PLUGIN_ID
}

func (api *Anthropic) GetPluginName() string {
	return PLUGIN_NAME
}

var AiApiPlugin Anthropic
