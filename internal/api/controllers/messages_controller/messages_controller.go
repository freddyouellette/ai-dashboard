package messages_controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/gorilla/mux"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type MessagesService interface {
	GetAll() ([]*models.Message, error)
	GetById(id uint) (*models.Message, error)
	Create(entity *models.Message) (*models.Message, error)
	GetChatMessages(chatId uint) ([]*models.Message, error)
}

type MessagesController struct {
	*entity_request_controller.EntityRequestController[models.Message]
	responseHandler ResponseHandler
	messagesService MessagesService
}

func NewMessagesController(
	entityRequestController *entity_request_controller.EntityRequestController[models.Message],
	responseHandler ResponseHandler,
	messagesService MessagesService,
) *MessagesController {
	return &MessagesController{
		EntityRequestController: entityRequestController,
		responseHandler:         responseHandler,
		messagesService:         messagesService,
	}
}

var (
	ErrInvalidId = errors.New("invalid id")
)

func (h *MessagesController) HandleGetMessageByChatIdRequest(w http.ResponseWriter, r *http.Request) {
	chatId, err := strconv.ParseUint(mux.Vars(r)["chat_id"], 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", ErrInvalidId, err.Error()))
		return
	}
	responseObject, err := h.messagesService.GetChatMessages(uint(chatId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}
