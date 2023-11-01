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
	HandleUpdateEntityByIdRequest(w http.ResponseWriter, r *http.Request)
	HandleDeleteEntityByIdRequest(w http.ResponseWriter, r *http.Request)
}

type ChatsController interface {
	EntityRequestController[models.Chat]
	HandleGetChatResponseRequest(w http.ResponseWriter, r *http.Request)
}

type MessagesController interface {
	EntityRequestController[models.Message]
	HandleGetMessageByChatIdRequest(w http.ResponseWriter, r *http.Request)
}

type RequestLogger interface {
	CreateRequestLoggerHandler(next http.Handler) http.Handler
}

func NewRouter(
	botsController EntityRequestController[models.Bot],
	chatsController ChatsController,
	messagesController MessagesController,
	requestLogger RequestLogger,
) http.Handler {
	router := mux.NewRouter()

	router.Use(requestLogger.CreateRequestLoggerHandler)

	router.HandleFunc("/api/bots", botsController.HandleGetAllEntitiesRequest).Methods("GET")
	router.HandleFunc("/api/bots", botsController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/api/bots", botsController.HandleUpdateEntityByIdRequest).Methods("PUT")
	router.HandleFunc("/api/bots/{id}", botsController.HandleGetEntityByIdRequest).Methods("GET")
	router.HandleFunc("/api/bots/{id}", botsController.HandleDeleteEntityByIdRequest).Methods("DELETE")

	router.HandleFunc("/api/chats", chatsController.HandleGetAllEntitiesRequest).Methods("GET")
	router.HandleFunc("/api/chats", chatsController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/api/chats/{id}", chatsController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/api/chats/{id}/response", chatsController.HandleGetChatResponseRequest).Methods("GET")

	router.HandleFunc("/api/messages", messagesController.HandleGetAllEntitiesRequest).Methods("GET")
	router.HandleFunc("/api/messages", messagesController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/api/messages/{id}", messagesController.HandleGetEntityByIdRequest).Methods("GET")
	router.HandleFunc("/api/chats/{chat_id}/messages", messagesController.HandleGetMessageByChatIdRequest).Methods("GET")

	return router
}
