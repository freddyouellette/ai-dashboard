package messages_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	responseHandler ResponseHandler
	messagesService MessagesService
}

func NewMessagesController(
	responseHandler ResponseHandler,
	messagesService MessagesService,
) *MessagesController {
	return &MessagesController{
		responseHandler: responseHandler,
		messagesService: messagesService,
	}
}

var (
	ErrInvalidId = errors.New("invalid id")
)

func (h *MessagesController) HandleGetAllMessagesRequest(w http.ResponseWriter, r *http.Request) {
	responseObject, err := h.messagesService.GetAll()
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *MessagesController) HandleCreateMessageRequest(w http.ResponseWriter, r *http.Request) {
	var message *models.Message
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, err)
		return
	}
	responseObject, err := h.messagesService.Create(message)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *MessagesController) HandleGetMessageByIdRequest(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", ErrInvalidId, err.Error()))
		return
	}
	responseObject, err := h.messagesService.GetById(uint(entityId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *MessagesController) HandleGetMessageByChatIdRequest(w http.ResponseWriter, r *http.Request) {
	chatId, err := strconv.ParseUint(mux.Vars(r)["chat_id"], 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", ErrInvalidId, err.Error()))
		return
	}
	responseObject, err := h.messagesService.GetChatMessages(uint(chatId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}
