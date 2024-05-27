package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/bots_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/chats_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/messages_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/error_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/logged_client"
	"github.com/freddyouellette/ai-dashboard/internal/api/request_logger"
	"github.com/freddyouellette/ai-dashboard/internal/api/response_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/router"
	"github.com/freddyouellette/ai-dashboard/internal/events/event_dispatcher"
	"github.com/freddyouellette/ai-dashboard/internal/events/event_handler"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/entity_repository"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/messages_repository"
	"github.com/freddyouellette/ai-dashboard/internal/repositories/users_repository"
	"github.com/freddyouellette/ai-dashboard/internal/services/chats_service"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
	"github.com/freddyouellette/ai-dashboard/internal/services/messages_service"
	"github.com/freddyouellette/ai-dashboard/internal/services/users_service"
	"github.com/freddyouellette/ai-dashboard/internal/util/logger"
	"github.com/freddyouellette/ai-dashboard/internal/util/plugin_loader"
	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	checkEnv()

	API_PORT := os.Getenv("API_PORT")
	frontendStr, ok := os.LookupEnv("FRONTEND")
	if !ok {
		frontendStr = "true"
	}
	FRONTEND, err := strconv.ParseBool(frontendStr)
	if err != nil {
		panic(err)
	}

	makeDir("data")

	db, err := gorm.Open(sqlite.Open("data/data.db"))
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Bot{})
	db.AutoMigrate(&models.Chat{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.User{})

	ERROR_LOG := os.Getenv("ERROR_LOG")

	errorFile, err := os.OpenFile(ERROR_LOG, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger(errorFile, logrus.ErrorLevel)
	errorHandler := error_handler.NewErrorHandler(logger)
	responseHandler := response_handler.NewResponseHandler(errorHandler)

	eventDispatcher := event_dispatcher.NewEventDispatcher(logger)

	httpClient := logged_client.NewLoggedClient(http.DefaultClient, logger, logged_client.Options{
		LogRequestHeaders:  true,
		LogRequestBody:     true,
		LogResponseHeaders: true,
		LogResponseBody:    true,
		PrettyJson:         true,
	})
	aiApis, err := plugin_loader.LoadPlugins[plugin_models.AiApiPlugin](os.Getenv("AI_API_PLUGINS"), "AiApiPlugin")
	if err != nil {
		panic(err)
	}

	for _, aiApi := range aiApis {
		aiApi.Initialize(&plugin_models.AiApiPluginOptions{
			Client: httpClient,
			Logger: logger,
		})
	}

	botsRepository := entity_repository.NewRepository[models.Bot](db)
	botsService := entity_service.NewEntityService[models.Bot](botsRepository)
	botsController := bots_controller.NewBotsController(
		entity_request_controller.NewEntityRequestController[models.Bot](
			responseHandler,
			botsService,
		),
		responseHandler,
		aiApis,
	)
	messagesRepository := messages_repository.NewMessagesRepository(
		entity_repository.NewRepository[models.Message](db),
		db,
	)
	messagesService := messages_service.NewMessagesService(
		entity_service.NewEntityService[models.Message](messagesRepository),
		messagesRepository,
		eventDispatcher,
	)
	chatsRepository := entity_repository.NewRepository[models.Chat](db)
	chatsService := chats_service.NewChatsService(
		entity_service.NewEntityService[models.Chat](chatsRepository),
		botsService,
		messagesService,
		aiApis,
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

	eventHandler := event_handler.NewEventHandler(chatsService, botsService)
	eventDispatcher.Register(models.EVENT_TYPE_MESSAGE_CREATED, eventHandler.HandleMessageCreatedEvent)

	usersRepository := users_repository.NewUsersRepository(entity_repository.NewRepository[models.User](db), db)
	usersService := users_service.NewUsersService(entity_service.NewEntityService[models.User](usersRepository), usersRepository)

	apiMiddlewares, err := plugin_loader.LoadPlugins[plugin_models.ApiMiddlewareFactory](os.Getenv("API_MIDDLEWARE_PLUGINS"), "ApiMiddleware")
	if err != nil {
		panic(err)
	}

	middlewareFuncs := make(map[string]func(http.Handler) http.Handler)
	for _, apiMiddleware := range apiMiddlewares {
		apiMiddleware.Initialize(&plugin_models.ApiMiddlewareOptions{
			UsersService: usersService,
		})
		middlewareFuncs[apiMiddleware.GetPluginId()] = apiMiddleware.Create
	}

	requestLogger := request_logger.NewRequestLogger(logger, request_logger.Options{
		LogHeaders:      false,
		LogRequestBody:  true,
		LogResponseBody: true,
		PrettyJson:      false,
	})
	apiRouter := router.NewRouter(FRONTEND, middlewareFuncs, botsController, chatsController, messagesController, requestLogger)

	apiRouter = cors.AllowAll().Handler(apiRouter)
	// router = cors.Default().Handler(router)

	fmt.Println("API listening on port " + API_PORT)
	err = http.ListenAndServe(":"+API_PORT, apiRouter)
	if err != nil {
		panic(err)
	}
}

func makeDir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
}

func checkEnv() {
	myEnv, err := godotenv.Read(".env.dist")
	if err != nil {
		panic("Error loading .env.dist file")
	}
	for key := range myEnv {
		if _, ok := os.LookupEnv(key); !ok {
			panic("Missing .env var " + key)
		}
	}
}
