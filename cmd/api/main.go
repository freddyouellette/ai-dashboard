package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/freddyouellette/ai-dashboard/internal/api/controllers/entity_request_controller"
	"github.com/freddyouellette/ai-dashboard/internal/api/error_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/response_handler"
	"github.com/freddyouellette/ai-dashboard/internal/api/router"
	"github.com/freddyouellette/ai-dashboard/internal/models"
	"github.com/freddyouellette/ai-dashboard/internal/repositories"
	"github.com/freddyouellette/ai-dashboard/internal/services/entity_service"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	API_PORT := os.Getenv("API_PORT")

	db, err := gorm.Open(sqlite.Open("data/data.db"))
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Bot{})

	errorHandler := error_handler.NewErrorHandler()
	responseHandler := response_handler.NewResponseHandler(errorHandler)
	botRepository := repositories.NewRepository[models.Bot](db)
	botService := entity_service.NewEntityService[models.Bot](botRepository)
	requestController := entity_request_controller.NewEntityRequestController[models.Bot](
		responseHandler,
		botService,
	)
	router := router.NewRouter(requestController)

	router = cors.AllowAll().Handler(router)
	// router = cors.Default().Handler(router)

	http.Handle("/", router)
	fmt.Println("Listening on port " + API_PORT)
	http.ListenAndServe(":"+API_PORT, nil)
}
