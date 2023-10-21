package entity_request_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ResponseHandler interface {
	HandleResponseObject(w http.ResponseWriter, response interface{}, err error)
}

type EntityService[e any] interface {
	GetAll() ([]*e, error)
	GetById(id uint) (*e, error)
	Create(entity *e) (*e, error)
	Update(entity *e) (*e, error)
}

type EntityRequestController[e any] struct {
	responseHandler ResponseHandler
	entityService   EntityService[e]
}

func NewEntityRequestController[e any](responseHandler ResponseHandler, entityService EntityService[e]) *EntityRequestController[e] {
	return &EntityRequestController[e]{
		responseHandler: responseHandler,
		entityService:   entityService,
	}
}

var (
	ErrInvalidId = errors.New("invalid id")
)

func (h *EntityRequestController[e]) HandleGetAllEntitiesRequest(w http.ResponseWriter, r *http.Request) {
	responseObject, err := h.entityService.GetAll()
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleCreateEntityRequest(w http.ResponseWriter, r *http.Request) {
	var entity *e
	err := json.NewDecoder(r.Body).Decode(entity)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, err)
		return
	}
	responseObject, err := h.entityService.Create(entity)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleGetEntityByIdRequest(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, fmt.Errorf("%w: %s", ErrInvalidId, err.Error()))
		return
	}
	responseObject, err := h.entityService.GetById(uint(entityId))
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}

func (h *EntityRequestController[e]) HandleUpdateEntityByIdRequest(w http.ResponseWriter, r *http.Request) {
	var entity *e
	err := json.NewDecoder(r.Body).Decode(entity)
	if err != nil {
		h.responseHandler.HandleResponseObject(w, nil, err)
		return
	}
	responseObject, err := h.entityService.Update(entity)
	h.responseHandler.HandleResponseObject(w, responseObject, err)
}
