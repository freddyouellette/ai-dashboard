package error_handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) HandleError(w http.ResponseWriter, err error) {
	var status int
	var message string
	switch {
	case errors.Is(err, models.ErrResourceNotFound):
		status = http.StatusNotFound
		message = "Error: " + err.Error()
	case errors.Is(err, models.ErrInvalidResourceSyntax):
		status = http.StatusBadRequest
		message = "Error: " + err.Error()
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
	}
	fmt.Printf("[%s] %s: %v\n", time.Now().Format(time.RFC1123Z), message, err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}
