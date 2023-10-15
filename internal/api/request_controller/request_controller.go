package request_controller

import (
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

type BotService interface {
	GetBots() ([]models.Bot, error)
	GetBotById(id uint) (*models.Bot, error)
}

type RequestController struct {
	responseHandler ResponseHandler
	botService      BotService
}

func NewRequestController(responseHandler ResponseHandler, botService BotService) *RequestController {
	return &RequestController{
		responseHandler: responseHandler,
		botService:      botService,
	}
}

var (
	ErrInvalidId = errors.New("invalid id")
)

func (h *RequestController) HandleGetBotsRequest(w http.ResponseWriter, r *http.Request) {
	responseObject, err := h.botService.GetBots()
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *RequestController) HandleGetBotRequest(w http.ResponseWriter, r *http.Request) {
	botID, err := strconv.ParseUint(mux.Vars(r)["bot_id"], 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", ErrInvalidId, err.Error()))
		return
	}
	responseObject, err := h.botService.GetBotById(uint(botID))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}
