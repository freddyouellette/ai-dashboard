package router

import (
	"github.com/freddyouellette/ai-dashboard/internal/api/request_controller"
	"github.com/gorilla/mux"
)

func NewRouter(controller *request_controller.RequestController) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/bots", controller.HandleGetBotsRequest).Methods("GET")
	router.HandleFunc("/bots/{bot_id}", controller.HandleGetBotRequest).Methods("GET")

	return router
}
