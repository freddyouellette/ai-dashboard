package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

const (
	PLUGIN_ID   = "google"
	PLUGIN_NAME = "Google"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Google struct {
	client       HttpClient
	googleApiKey string
}

type googleMessage struct {
	Text string `json:"text"`
}

type googleContents struct {
	Parts []googleMessage `json:"parts"`
	Role  string          `json:"role,omitempty"`
}

type googleGenerateContentRequest struct {
	SystemInstruction []googleContents `json:"systemInstruction,omitempty"`
	Contents          []googleContents `json:"contents"`
	GenerationConfig  struct {
		Temperature float64 `json:"temperature,omitempty"`
	} `json:"generationConfig,omitempty"`
}

type googleChatCompletionResponse struct {
	Candidates []struct {
		Content struct {
			Parts []googleMessage `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type googleModelsResponse struct {
	Models []struct {
		Name                       string   `json:"name"`
		SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
		DisplayName                string   `json:"displayName"`
		Description                string   `json:"description"`
		Temperature                float64  `json:"temperature"`
		MaxTemperature             *float64 `json:"maxTemperature,omitempty"`
	} `json:"models"`
}

var (
	ErrBadResponse = errors.New("received bad response")
)

func (api *Google) Initialize(options *plugin_models.AiApiPluginOptions) error {
	api.client = options.Client
	api.googleApiKey = os.Getenv("GOOGLE_API_KEY")
	return nil
}

func (api *Google) CompleteChat(chatCompletionRequest *plugin_models.ChatCompletionRequest) (*plugin_models.ChatCompletionResponse, error) {
	requestBody := googleGenerateContentRequest{}

	systemInstructionMessages := make([]googleContents, 0)
	requestMessages := make([]googleContents, 0)
	for _, message := range chatCompletionRequest.Messages {
		var messageRole string
		switch message.Role {
		case plugin_models.ChatCompletionRoleUser:
			messageRole = "user"
		case plugin_models.ChatCompletionRoleAssistant:
			messageRole = "model"
		case plugin_models.ChatCompletionRoleSystem:
			messageRole = "model"
			systemInstructionMessages = append(systemInstructionMessages, googleContents{
				Role: messageRole,
				Parts: []googleMessage{
					{
						Text: message.Content,
					},
				},
			})
			continue
		default:
			continue
		}

		requestMessages = append(requestMessages, googleContents{
			Role: messageRole,
			Parts: []googleMessage{
				{
					Text: message.Content,
				},
			},
		})
	}

	if len(systemInstructionMessages) > 0 {
		requestBody.SystemInstruction = append(requestBody.SystemInstruction, systemInstructionMessages...)
	}
	if len(requestMessages) > 0 {
		requestBody.Contents = append(requestBody.Contents, requestMessages...)
	}
	// google temperature goes from 0 to 2
	requestBody.GenerationConfig.Temperature = chatCompletionRequest.Temperature * 2

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	requestReader := bytes.NewBuffer(jsonBody)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s", chatCompletionRequest.Model, api.googleApiKey)

	request, err := http.NewRequest(http.MethodPost, url, requestReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody googleChatCompletionResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %s", ErrBadResponse, response.StatusCode, responseBody)
	}

	return &plugin_models.ChatCompletionResponse{
		Message: &plugin_models.ChatCompletionMessage{
			Content: responseBody.Candidates[0].Content.Parts[0].Text,
			Role:    plugin_models.ChatCompletionRoleAssistant,
		},
	}, nil
}

func (api *Google) GetModels() (*plugin_models.GetModelsResponse, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", api.googleApiKey)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody googleModelsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: [%d] %v", ErrBadResponse, response.StatusCode, responseBody)
	}

	botModels := make([]*plugin_models.AiModel, 0)

	requiredGenerationMethods := []string{"generateContent"}
	for _, model := range responseBody.Models {
		hasRequiredMethod := false
		for _, method := range model.SupportedGenerationMethods {
			for _, requiredMethod := range requiredGenerationMethods {
				if method == requiredMethod {
					hasRequiredMethod = true
					break
				}
			}
			if hasRequiredMethod {
				break
			}
		}
		if !hasRequiredMethod {
			continue
		}
		temp := model.Temperature
		if model.MaxTemperature != nil {
			temp /= *model.MaxTemperature
		}
		botModels = append(botModels, &plugin_models.AiModel{
			Id:         model.Name,
			Name:       model.DisplayName,
			AuthorId:   PLUGIN_ID,
			AuthorName: PLUGIN_NAME,
		})
	}

	return &plugin_models.GetModelsResponse{
		Models: botModels,
	}, nil
}

func (api *Google) GetPluginId() string {
	return PLUGIN_ID
}

func (api *Google) GetPluginName() string {
	return PLUGIN_NAME
}

var AiApiPlugin Google
