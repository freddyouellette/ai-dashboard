package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

const (
	PLUGIN_ID   = "openai"
	PLUGIN_NAME = "OpenAI"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OpenAi struct {
	client              HttpClient
	maxCompletionTokens int
	chatGptAccessToken  string
}

type openAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAiChatCompletionRequest struct {
	Model               string          `json:"model"`
	Messages            []openAiMessage `json:"messages"`
	MaxCompletionTokens int             `json:"max_completion_tokens"`
	Temperature         *float64        `json:"temperature,omitempty"`
}

type openAiChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAiModel struct {
	Id      string `json:"id"`
	Created int64  `json:"created"`
}

type openAiModelsResponse struct {
	Data []openAiModel `json:"data"`
}

var (
	ErrBadResponse = errors.New("received bad response")
)

func (api *OpenAi) Initialize(options *plugin_models.AiApiPluginOptions) error {
	api.client = options.Client
	maxCompletionTokens, err := strconv.Atoi(os.Getenv("OPENAI_MAX_COMPLETION_TOKENS"))
	if err != nil {
		return err
	}
	api.maxCompletionTokens = maxCompletionTokens
	api.chatGptAccessToken = os.Getenv("OPENAI_ACCESS_TOKEN")
	return nil
}

func (api *OpenAi) CompleteChat(chatCompletionRequest *plugin_models.ChatCompletionRequest) (*plugin_models.ChatCompletionResponse, error) {
	modelRegex := regexp.MustCompile(`^o[0-9]+-`)
	// o1 family doesn't use system prompts, or temperature, for some reason
	supportsSystemRole := modelRegex.MatchString(chatCompletionRequest.Model)
	supportsTemperature := !modelRegex.MatchString(chatCompletionRequest.Model)

	requestBody := openAiChatCompletionRequest{
		Model:               chatCompletionRequest.Model,
		MaxCompletionTokens: api.maxCompletionTokens,
	}

	if supportsTemperature {
		t := chatCompletionRequest.Temperature * 2
		requestBody.Temperature = &t
	}

	requestMessages := make([]openAiMessage, 0)
	for _, message := range chatCompletionRequest.Messages {
		var messageRole string
		switch message.Role {
		case plugin_models.ChatCompletionRoleUser:
			messageRole = "user"
		case plugin_models.ChatCompletionRoleAssistant:
			messageRole = "assistant"
		case plugin_models.ChatCompletionRoleSystem:
			if supportsSystemRole {
				messageRole = "user"
			} else {
				messageRole = "developer"
			}
		default:
			continue
		}

		requestMessages = append(requestMessages, openAiMessage{
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

	var responseBody openAiChatCompletionResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %s", ErrBadResponse, response.StatusCode, responseBody)
	}

	return &plugin_models.ChatCompletionResponse{
		Message: &plugin_models.ChatCompletionMessage{
			Content: responseBody.Choices[0].Message.Content,
			Role:    plugin_models.ChatCompletionRoleAssistant,
		},
	}, nil
}

func (api *OpenAi) GetModels() (*plugin_models.GetModelsResponse, error) {
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

	var responseBody openAiModelsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %v", ErrBadResponse, response.StatusCode, responseBody)
	}

	botModels := make([]*plugin_models.AiModel, 0)

	for _, model := range responseBody.Data {
		botModels = append(botModels, &plugin_models.AiModel{
			Id:         model.Id,
			AuthorId:   PLUGIN_ID,
			AuthorName: PLUGIN_NAME,
			CreatedAt:  time.Unix(model.Created, 0),
		})
	}

	return &plugin_models.GetModelsResponse{
		Models: botModels,
	}, nil
}

func (api *OpenAi) GetPluginId() string {
	return PLUGIN_ID
}

func (api *OpenAi) GetPluginName() string {
	return PLUGIN_NAME
}

var AiApiPlugin OpenAi
