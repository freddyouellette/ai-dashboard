package bots_controller

import (
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/user_scoped_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type BotsController struct {
	*user_scoped_request_controller.UserScopedRequestController[*models.Bot]
	responseHandler ResponseHandler
	aiApis          map[string]plugin_models.AiApiPlugin
}

func NewBotsController(
	userScopedRequestController *user_scoped_request_controller.UserScopedRequestController[*models.Bot],
	responseHandler ResponseHandler,
	aiApis map[string]plugin_models.AiApiPlugin,
) *BotsController {
	return &BotsController{
		UserScopedRequestController: userScopedRequestController,
		responseHandler:             responseHandler,
		aiApis:                      aiApis,
	}
}

func (h *BotsController) HandleGetBotModelsRequest(w http.ResponseWriter, r *http.Request) {
	botModels := make([]*plugin_models.AiModel, 0)
	for _, aiApi := range h.aiApis {
		modelsResponse, err := aiApi.GetModels()
		if err != nil {
			h.responseHandler.HandleResponseObject(w, nil, err)
			return
		}
		botModels = append(botModels, modelsResponse.Models...)
	}
	h.responseHandler.HandleResponseObject(w, botModels, nil)
}
