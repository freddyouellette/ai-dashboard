package main

import (
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/api/error_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/response_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/router"
	"github.com/freddyouellette/ai-dashboard/internal/services/bots"
)

func main() {
	errorHandler := error_handler.NewErrorHandler()
	responseHandler := response_handler.NewResponseHandler(errorHandler)
	botService := bots.NewBotService()
	requestController := request_controller.NewRequestController(
		responseHandler,
		botService,
	)
	router := router.NewRouter(requestController)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
