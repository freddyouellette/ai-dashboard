package chats_controller

import (
	"net/http"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/user_scoped_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
	"github.com/go-chi/chi/v5"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type ChatsService interface {
	GetChatResponse(userId uint, chatId uint) (*models.Message, error)
	GetMessageCorrection(messageId uint) (*models.Message, error)
}

type RequestUtils interface {
	GetContextInt(r *http.Request, key any, def int) int
}

type ChatsController struct {
	*user_scoped_request_controller.UserScopedRequestController[*models.Chat]
	responseHandler ResponseHandler
	chatsService    ChatsService
	requestUtils    RequestUtils
}

func NewChatsController(
	userScopedRequestController *user_scoped_request_controller.UserScopedRequestController[*models.Chat],
	responseHandler ResponseHandler,
	chatsService ChatsService,
	requestUtils RequestUtils,
) *ChatsController {
	return &ChatsController{
		UserScopedRequestController: userScopedRequestController,
		responseHandler:             responseHandler,
		chatsService:                chatsService,
		requestUtils:                requestUtils,
	}
}

func (h *ChatsController) HandleGetChatResponseRequest(w http.ResponseWriter, r *http.Request) {
	chatId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, err)
		return
	}

	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	responseObject, err := h.chatsService.GetChatResponse(uint(userId), uint(chatId))

	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *ChatsController) HandleGetMessageCorrectionRequest(w http.ResponseWriter, r *http.Request) {
	messageId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, err)
		return
	}

	message, err := h.chatsService.GetMessageCorrection(uint(messageId))

	h.responseHandler.HandleResponseObject(w, message, err)
}
