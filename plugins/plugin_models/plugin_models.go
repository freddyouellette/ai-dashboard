package plugin_models

import (
	"net/http"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/api/logged_client"
	"github.com/freddyouellette/ai-dashboard/internal/services/users_service"
	"github.com/freddyouellette/ai-dashboard/internal/util/logger"
)

type Plugin interface {
	GetPluginId() string
}

type ChatCompletionRole string

const (
	ChatCompletionRoleUser      ChatCompletionRole = "user"
	ChatCompletionRoleAssistant ChatCompletionRole = "assistant"
	ChatCompletionRoleSystem    ChatCompletionRole = "system"
)

type ChatCompletionMessage struct {
	Content string             `json:"content"`
	Role    ChatCompletionRole `json:"role"`
}

type ChatCompletionRequest struct {
	Model    string                   `json:"model"`
	Messages []*ChatCompletionMessage `json:"messages"`
	// Temperature is a float value between 0 and 1. Higher values will result in more creative completions.
	Temperature float64 `json:"temperature"`
}

type ChatCompletionResponse struct {
	Message *ChatCompletionMessage `json:"message"`
}

type GetModelsResponse struct {
	Models []*AiModel `json:"models"`
}

type AiModel struct {
	Id string `json:"id"`
	// Person or Organization that created the model, e.g. "OpenAI"
	AuthorId   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type AiApiPluginOptions struct {
	Client *logged_client.LoggedClient
	Logger *logger.Logger
}

type AiApiPlugin interface {
	Initialize(options *AiApiPluginOptions) error
	GetPluginId() string
	GetPluginName() string
	CompleteChat(chatCompletionRequest *ChatCompletionRequest) (*ChatCompletionResponse, error)
	GetModels() (*GetModelsResponse, error)
}

type UserIdContextKey struct{}

type ApiMiddlewareOptions struct {
	UsersService *users_service.UsersService
}

type ApiMiddlewareFactory interface {
	Initialize(options *ApiMiddlewareOptions) error
	GetPluginId() string
	Create(next http.Handler) http.Handler
}
