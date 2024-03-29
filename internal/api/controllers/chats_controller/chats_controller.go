package chats_controller

import (
	"net/http"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/go-chi/chi/v5"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type ChatsService interface {
	GetChatResponse(chatId uint) (*models.Message, error)
	GetMessageCorrection(messageId uint) (*models.Message, error)
}

type ChatsController struct {
	*entity_request_controller.EntityRequestController[models.Chat]
	responseHandler ResponseHandler
	chatsService    ChatsService
}

func NewChatsController(
	entityRequestController *entity_request_controller.EntityRequestController[models.Chat],
	responseHandler ResponseHandler,
	chatsService ChatsService,
) *ChatsController {
	return &ChatsController{
		EntityRequestController: entityRequestController,
		responseHandler:         responseHandler,
		chatsService:            chatsService,
	}
}

func (h *ChatsController) HandleGetChatResponseRequest(w http.ResponseWriter, r *http.Request) {
	chatId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, err)
		return
	}

	responseObject, err := h.chatsService.GetChatResponse(uint(chatId))

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
