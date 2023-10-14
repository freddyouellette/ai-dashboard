package request_controller

import (
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/gorilla/mux"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type BotService interface {
	GetBots() ([]models.Bot, error)
	GetBotById(id string) (*models.Bot, error)
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

func (h *RequestController) HandleGetBotsRequest(w http.ResponseWriter, r *http.Request) {
	responseObject, err := h.botService.GetBots()
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *RequestController) HandleGetBotRequest(w http.ResponseWriter, r *http.Request) {
	responseObject, err := h.botService.GetBotById(mux.Vars(r)["bot_uuid"])
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}
