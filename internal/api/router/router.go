package router

import (
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/gorilla/mux"
)

type EntityRequestController[e any] interface {
	HandleGetAllEntitiesRequest(w http.ResponseWriter, r *http.Request)
	HandleGetEntityByIdRequest(w http.ResponseWriter, r *http.Request)
	HandleCreateEntityRequest(w http.ResponseWriter, r *http.Request)
}

func NewRouter(botController EntityRequestController[models.Bot]) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/bots", botController.HandleGetAllEntitiesRequest).Methods("GET")
	router.HandleFunc("/bots", botController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/bots/{bot_id}", botController.HandleGetEntityByIdRequest).Methods("GET")

	return router
}
