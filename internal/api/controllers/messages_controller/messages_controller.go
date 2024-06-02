package messages_controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/user_scoped_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

const (
	DEFAULT_PAGE_SIZE = 20
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type MessagesService interface {
	GetAllPaginated(userId uint, options *models.GetMessagesOptions) (*models.MessagesDTO, error)
	CreateWithUserId(entity *models.Message, userId uint) (*models.Message, error)
}

type RequestUtils interface {
	GetQueryInt(r *http.Request, param string, def int) (int, error)
	GetContextInt(r *http.Request, key any, def int) int
}

type MessagesController struct {
	*user_scoped_request_controller.UserScopedRequestController[*models.Message]
	responseHandler ResponseHandler
	messagesService MessagesService
	requestUtils    RequestUtils
}

func NewMessagesController(
	userScopedRequestController *user_scoped_request_controller.UserScopedRequestController[*models.Message],
	responseHandler ResponseHandler,
	messagesService MessagesService,
	requestUtils RequestUtils,
) *MessagesController {
	return &MessagesController{
		UserScopedRequestController: userScopedRequestController,
		responseHandler:             responseHandler,
		messagesService:             messagesService,
		requestUtils:                requestUtils,
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
	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	message.Role = models.MESSAGE_ROLE_USER
	responseObject, err := h.messagesService.CreateWithUserId(&message, uint(userId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *MessagesController) HandleGetAllPaginatedRequest(w http.ResponseWriter, r *http.Request) {
	options := &models.GetMessagesOptions{}
	chatID, err := h.requestUtils.GetQueryInt(r, "chat_id", 0)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	options.ChatID = uint(chatID)
	page, err := h.requestUtils.GetQueryInt(r, "page", 1)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	options.Page = page
	perPage, err := h.requestUtils.GetQueryInt(r, "per_page", DEFAULT_PAGE_SIZE)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, models.ErrInvalidResourceSyntax)
		return
	}
	options.PerPage = perPage
	userId := h.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
	responseObject, err := h.messagesService.GetAllPaginated(uint(userId), options)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}
