package bots_controller

import (
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type AiApi interface {
	GetBotModels() ([]*models.BotModel, error)
}

type BotsController struct {
	*entity_request_controller.EntityRequestController[models.Bot]
	responseHandler ResponseHandler
	aiApi           AiApi
}

func NewBotsController(
	entityRequestController *entity_request_controller.EntityRequestController[models.Bot],
	responseHandler ResponseHandler,
	aiApi AiApi,
) *BotsController {
	return &BotsController{
		EntityRequestController: entityRequestController,
		responseHandler:         responseHandler,
		aiApi:                   aiApi,
	}
}

func (h *BotsController) HandleGetBotModelsRequest(w http.ResponseWriter, r *http.Request) {
	botModels, err := h.aiApi.GetBotModels()
	h.responseHandler.HandleResponseObject(w, botModels, err)
}
