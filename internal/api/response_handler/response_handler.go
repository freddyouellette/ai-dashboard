package response_handler

import (
	"encoding/json"
	"net/http"
)

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, err error)
}

type ResponseHandler struct {
	errorHandler ErrorHandler
}

func NewResponseHandler(errorHandler ErrorHandler) *ResponseHandler {
	return &ResponseHandler{
		errorHandler: errorHandler,
	}
}

func (h *ResponseHandler) HandleResponseObject(w http.ResponseWriter, response interface{}, err error) {
	if err != nil {
		h.errorHandler.HandleError(w, err)
		return
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		h.errorHandler.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}
