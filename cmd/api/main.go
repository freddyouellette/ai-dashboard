package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/chats_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/messages_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/error_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/logged_client"
	"github.com/freddyouellette/ai-dashboard/internal/api/request_logger"
	"github.com/freddyouellette/ai-dashboard/internal/api/response_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/router"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/entity_repository"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/messages_repository"
	"github.com/freddyouellette/ai-dashboard/internal/services/ai_api.go"
	"github.com/freddyouellette/ai-dashboard/internal/services/chats_service"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
	"github.com/freddyouellette/ai-dashboard/internal/services/messages_service"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	API_PORT := os.Getenv("API_PORT")
	OPENAI_ACCESS_TOKEN := os.Getenv("OPENAI_ACCESS_TOKEN")

	db, err := gorm.Open(sqlite.Open("data/data.db"))
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Bot{})
	db.AutoMigrate(&models.Chat{})
	db.AutoMigrate(&models.Message{})

	errorHandler := error_handler.NewErrorHandler()
	responseHandler := response_handler.NewResponseHandler(errorHandler)
	botsRepository := entity_repository.NewRepository[models.Bot](db)
	botsService := entity_service.NewEntityService[models.Bot](botsRepository)
	botsController := entity_request_controller.NewEntityRequestController[models.Bot](
		responseHandler,
		botsService,
	)
	messagesRepository := messages_repository.NewMessagesRepository(
		entity_repository.NewRepository[models.Message](db),
		db,
	)
	messagesService := messages_service.NewMessagesService(
		entity_service.NewEntityService[models.Message](messagesRepository),
		messagesRepository,
	)
	chatsRepository := entity_repository.NewRepository[models.Chat](db)
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	httpClient := logged_client.NewLoggedClient(http.DefaultClient, logger, logged_client.Options{
		LogRequestHeaders:  true,
		LogRequestBody:     true,
		LogResponseHeaders: true,
		LogResponseBody:    true,
		PrettyJson:         true,
	})
	// httpClient := http.DefaultClient
	aiApi := ai_api.NewAiApi(httpClient, 5000, "https://api.openai.com/v1/chat/completions", OPENAI_ACCESS_TOKEN)
	chatsService := chats_service.NewChatsService(
		entity_service.NewEntityService[models.Chat](chatsRepository),
		botsService,
		messagesService,
		aiApi,
	)
	chatsController := chats_controller.NewChatsController(
		entity_request_controller.NewEntityRequestController[models.Chat](
			responseHandler,
			chatsService,
		),
		responseHandler,
		chatsService,
	)
	messagesController := messages_controller.NewMessagesController(
		entity_request_controller.NewEntityRequestController[models.Message](
			responseHandler,
			messagesService,
		),
		responseHandler,
		messagesService,
	)
	requestLogger := request_logger.NewRequestLogger(logger, request_logger.Options{
		LogHeaders:      false,
		LogRequestBody:  true,
		LogResponseBody: true,
		PrettyJson:      true,
	})
	apiRouter := router.NewRouter(botsController, chatsController, messagesController, requestLogger)

	apiRouter = cors.AllowAll().Handler(apiRouter)
	// router = cors.Default().Handler(router)

	// frontend
	WEB_PORT := os.Getenv("WEB_PORT")
	fs := http.FileServer(http.Dir("web/build"))
	frontendServer := http.NewServeMux()
	frontendServer.Handle("/", fs)
	fmt.Println("Frontend listening on port " + WEB_PORT)

	go http.ListenAndServe(":"+WEB_PORT, frontendServer)

	fmt.Println("API listening on port " + API_PORT)
	http.ListenAndServe(":"+API_PORT, apiRouter)
}
