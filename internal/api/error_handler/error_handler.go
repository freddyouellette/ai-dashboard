package error_handler

import (
	"errors"
	"net/http"

	"github.com/freddyouellette/ai-dashboard/internal/models"
)

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Warning(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
}

type ErrorHandler struct {
	logger Logger
}

func NewErrorHandler(logger Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) HandleError(w http.ResponseWriter, err error) {
	var status int
	var message string
	switch {
	case errors.Is(err, models.ErrResourceNotFound):
		status = http.StatusNotFound
		message = "Error: " + err.Error()
		h.logger.Warning("Resource not found", map[string]interface{}{"error": err.Error()})
	case errors.Is(err, models.ErrInvalidResourceSyntax):
		status = http.StatusBadRequest
		message = "Error: " + err.Error()
		h.logger.Warning("Invalid Resource Syntax", map[string]interface{}{"error": err.Error()})
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
		h.logger.Error("Internal Server Error", map[string]interface{}{"error": err.Error()})
	}
	w.WriteHeader(status)
	w.Write([]byte(message))
}
