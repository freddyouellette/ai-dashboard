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
}

type ChatsController interface {
	EntityRequestController[models.Chat]
	HandleGetChatResponseRequest(w http.ResponseWriter, r *http.Request)
}

type MessagesController interface {
	HandleGetAllMessagesRequest(w http.ResponseWriter, r *http.Request)
	HandleCreateMessageRequest(w http.ResponseWriter, r *http.Request)
	HandleGetMessageByIdRequest(w http.ResponseWriter, r *http.Request)
	HandleGetMessageByChatIdRequest(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	botsController EntityRequestController[models.Bot],
	chatsController ChatsController,
	messagesController MessagesController,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/bots", botsController.HandleGetAllEntitiesRequest).Methods("GET")
	router.HandleFunc("/bots", botsController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/bots", botsController.HandleUpdateEntityByIdRequest).Methods("PUT")
	router.HandleFunc("/bots/{id}", botsController.HandleGetEntityByIdRequest).Methods("GET")

	router.HandleFunc("/chats", chatsController.HandleGetAllEntitiesRequest).Methods("GET")
	router.HandleFunc("/chats", chatsController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/chats/{id}", chatsController.HandleCreateEntityRequest).Methods("POST")
	router.HandleFunc("/chats/{id}/response", chatsController.HandleGetChatResponseRequest).Methods("GET")

	router.HandleFunc("/messages", messagesController.HandleGetAllMessagesRequest).Methods("GET")
	router.HandleFunc("/messages", messagesController.HandleCreateMessageRequest).Methods("POST")
	router.HandleFunc("/messages/{id}", messagesController.HandleGetMessageByIdRequest).Methods("GET")
	router.HandleFunc("/chats/{chat_id}/messages", messagesController.HandleGetMessageByChatIdRequest).Methods("GET")

	return router
}
