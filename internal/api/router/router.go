package router

import (
	"net/http"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/api/web_handler"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type EntityRequestController[e any] interface {
	HandleGetAllEntitiesRequest(w http.ResponseWriter, r *http.Request)
	HandleGetEntityByIdRequest(w http.ResponseWriter, r *http.Request)
	HandleCreateEntityRequest(w http.ResponseWriter, r *http.Request)
	HandleUpdateEntityRequest(w http.ResponseWriter, r *http.Request)
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
	frontend bool,
	botsController EntityRequestController[models.Bot],
	chatsController ChatsController,
	messagesController MessagesController,
	requestLogger RequestLogger,
) http.Handler {
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{fmt.Sprintf("%s:%s", WEB_HOST, WEB_PORT)},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(30 * time.Second))

	router.Get("/api/bots", botsController.HandleGetAllEntitiesRequest)
	router.Post("/api/bots", botsController.HandleCreateEntityRequest)
	router.Put("/api/bots", botsController.HandleUpdateEntityRequest)
	router.Get("/api/bots/{id}", botsController.HandleGetEntityByIdRequest)
	router.Delete("/api/bots/{id}", botsController.HandleDeleteEntityByIdRequest)

	router.Get("/api/chats", chatsController.HandleGetAllEntitiesRequest)
	router.Post("/api/chats", chatsController.HandleCreateEntityRequest)
	router.Put("/api/chats", chatsController.HandleUpdateEntityRequest)
	router.Post("/api/chats/{id}", chatsController.HandleCreateEntityRequest)
	router.Get("/api/chats/{id}/response", chatsController.HandleGetChatResponseRequest)

	router.Get("/api/messages", messagesController.HandleGetAllEntitiesRequest)
	router.Post("/api/messages", messagesController.HandleCreateEntityRequest)
	router.Get("/api/messages/{id}", messagesController.HandleGetEntityByIdRequest)
	router.Get("/api/chats/{chat_id}/messages", messagesController.HandleGetMessageByChatIdRequest)

	if frontend {
		webHandler := web_handler.NewWebHandler("web/build", "index.html")
		router.HandleFunc("/*", webHandler.Handler)
	}

	return router
}
