package messages_controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
)

const (
	DEFAULT_PAGE_SIZE = 2
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type MessagesService interface {
	GetAllPaginated(options *models.GetMessagesOptions) (*models.MessagesDTO, error)
	GetById(id uint) (*models.Message, error)
	Create(entity *models.Message) (*models.Message, error)
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

func (h *MessagesController) HandleCreateEntityRequest(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	message.Role = models.MESSAGE_ROLE_USER
	responseObject, err := h.messagesService.Create(&message)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *MessagesController) HandleGetAllPaginatedRequest(w http.ResponseWriter, r *http.Request) {
	options := &models.GetMessagesOptions{}
	chatID, err := h.parseQueryParamToInt(r, "chat_id", 0)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	options.ChatID = uint(chatID)
	page, err := h.parseQueryParamToInt(r, "page", 1)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	options.Page = page
	perPage, err := h.parseQueryParamToInt(r, "per_page", DEFAULT_PAGE_SIZE)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	options.PerPage = perPage
	responseObject, err := h.messagesService.GetAllPaginated(options)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *MessagesController) parseQueryParamToInt(r *http.Request, param string, def int) (int, error) {
	paramStr := r.URL.Query().Get(param)
	if paramStr == "" {
		return def, nil
	}
	paramInt, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	return paramInt, nil
}
